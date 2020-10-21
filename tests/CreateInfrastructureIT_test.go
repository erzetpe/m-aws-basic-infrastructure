package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/resourcegroups"
	"github.com/epiphany-platform/aws-basic-infrastructure/tests/utils"
)

const (
	imageTag   = "epiphanyplatform/awsbi:0.0.1"
	awsTag     = "awsbi-module"
	moduleName = "awsbi-module"
	sshKeyName = "vms_rsa"
	path       = "./"
)

var (
	awsAccessKey              = "M_AWS_ACCESS_KEY=" + os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey              = "M_AWS_SECRET_KEY=" + os.Getenv("AWS_SECRET_ACCESS_KEY")
	sharedFilePath            = filepath.Join(path, "shared")
	sharedAbsoluteFilePath, _ = filepath.Abs(sharedFilePath)
	mountDir                  = sharedAbsoluteFilePath + ":/shared"
	dockerExecPath, _         = exec.LookPath("docker")
)

func TestMain(m *testing.M) {

	cleanup()
	setup()
	log.Println("Run tests")
	exitVal := m.Run()
	cleanup()
	log.Println("Finish test")
	os.Exit(exitVal)
}

func setup() {

	log.Println("Initialize test")
	if len(awsAccessKey) == 0 {
		log.Fatalf("expected non-empty AWS_ACCESS_KEY environment variable")
		os.Exit(1)
	}

	if len(awsSecretKey) == 0 {
		log.Fatalf("expected non-empty AWS_SECRET_KEY environment variable")
		os.Exit(1)
	}

	if _, err := os.Stat(sharedAbsoluteFilePath); os.IsNotExist(err) {
		os.Mkdir(sharedAbsoluteFilePath, os.ModePerm)
	}

	log.Println("Generating Keys")
	utils.GenerateRsaKeyPair(sharedAbsoluteFilePath, sshKeyName)
}

func cleanup() {
	log.Println("Starting cleanup.")
	err := os.RemoveAll(sharedAbsoluteFilePath)
	if err != nil {
		log.Fatal("Cannot remove data folder. ", err)
		os.Exit(1)
	}
	log.Println("Cleanup finished.")
}

func TestInitShouldCreateProperFileAndFolder(t *testing.T) {
	// given
	stateFilePath := "shared/state.yml"
	expectedFileContentRegexp := "kind: state\nawsbi:\n  status: initialized"

	// when
	utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "init", "M_NAME="+moduleName)

	data, err := ioutil.ReadFile(stateFilePath)

	if err != nil {
		t.Fatal("Cannot read state file: ", stateFilePath)
	}

	fileContent := string(data)

	matched, _ := regexp.MatchString(expectedFileContentRegexp, fileContent)

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedFileContentRegexp, "\nbut found:\n", fileContent)
	}

}

func TestOnPlanWithDefaultsShouldDisplayPlan(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer

	expectedOutputRegexp := ".*Plan: 14 to add, 0 to change, 0 to destroy.*"

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "plan", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, _ := regexp.MatchString(expectedOutputRegexp, outStr)

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}

}

func TestOnApplyShouldCreateEnvironment(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer
	expectedOutputRegexp := ".*Apply complete! Resources: 14 added, 0 changed, 0 destroyed.*"

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "apply", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, _ := regexp.MatchString(expectedOutputRegexp, outStr)

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}

}

func TestShouldCheckNumberOfVms(t *testing.T) {
	// given
	instancesNumber := 1

	session, err := session.NewSession(&aws.Config{Region: aws.String("eu-central-1")})
	if err != nil {
		t.Fatal("Cannot get session.")
	}

	// when

	// ec2
	ec2Client := ec2.New(session)

	filterParams := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(awsTag),
				},
			},
		},
	}

	ec2Result, err := ec2Client.DescribeInstances(filterParams)

	// rg
	rgClient := resourcegroups.New(session)

	rgName := "rg-" + moduleName

	rgResult, err := rgClient.GetGroup(&resourcegroups.GetGroupInput{
		GroupName: aws.String(rgName),
	})

	log.Println("rg", rgResult.Group)

	rgResourcesList, err := rgClient.ListGroupResources(&resourcegroups.ListGroupResourcesInput{
		GroupName: aws.String(rgName),
	})

	log.Println(rgResourcesList)

	// then
	if err != nil {
		t.Fatal("There was an error. ", err)
	}

	if len(ec2Result.Reservations[0].Instances) != 1 {
		t.Error("Expected ", instancesNumber, "instance, got ", len(ec2Result.Reservations[0].Instances))
	}

}

func TestOnDestroyPlanShouldDisplayDestroyPlan(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer
	expectedOutputRegexp := "Plan: 0 to add, 0 to change, 14 to destroy"

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "plan-destroy", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, _ := regexp.MatchString(expectedOutputRegexp, outStr)

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}
}

func TestOnDestroyShouldDestroyEnvironment(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer

	expectedOutputRegexp := "Apply complete! Resources: 0 added, 0 changed, 14 destroyed."

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "destroy", awsAccessKey, awsSecretKey)

	if stderr.Len() > 0 {
		t.Fatal("There was an error during executing a command. ", string(stderr.Bytes()))
	}

	outStr := string(stdout.Bytes())

	matched, _ := regexp.MatchString(expectedOutputRegexp, outStr)

	// then
	if !matched {
		t.Error("Expected to find expression matching:\n", expectedOutputRegexp, "\nbut found:\n", outStr)
	}
}
