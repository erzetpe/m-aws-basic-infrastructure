package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/epiphany-platform/aws-basic-infrastructure/tests/utils"
)

const (
	imageTag   = "epiphanyplatform/awsbi:0.0.1"
	moduleName = "awsbi-module"
	sshKeyName = "vms_rsa"
	path       = "./"
)

var (
	awsAccessKey              string = "M_AWS_ACCESS_KEY=" + os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey              string = "M_AWS_SECRET_KEY=" + os.Getenv("AWS_SECRET_KEY")
	sharedFilePath                   = filepath.Join(path, "shared")
	sharedAbsoluteFilePath, _        = filepath.Abs(sharedFilePath)
	mountDir                         = sharedAbsoluteFilePath + ":/shared"
	dockerExecPath, _                = exec.LookPath("docker")
)

func TestMain(m *testing.M) {

	setup()
	log.Println("Run tests")
	exitVal := m.Run()
	cleanup()

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

	log.Println("sharedAbsoluteFilePath: " + sharedAbsoluteFilePath)

	if something, err := os.Stat(sharedAbsoluteFilePath); os.IsNotExist(err) {
		fmt.Println(something)
		os.Mkdir(sharedAbsoluteFilePath, os.ModePerm)
	}

	log.Println("generateKeys")
	utils.GenerateRsaKeyPair(sharedAbsoluteFilePath, sshKeyName)
}

func cleanup() {
	log.Println("Cleanup.")
	os.RemoveAll(sharedAbsoluteFilePath)
	log.Println("Finish test")
}

func TestInitShouldCreateProperFileAndFolder(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "init", "M_NAME="+moduleName)

	// then
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)

	// if diff := deep.Equal(output, tt.wantOutput); diff != nil {
	// 	t.Error(diff)
	// }
}

func TestOnPlanShouldDisplayPlan(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "plan", awsAccessKey, awsSecretKey)

	// then
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)

}

func TestOnApplyShouldCreateEnvironment(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "apply", awsAccessKey, awsSecretKey)

	// then
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func TestOnDestroyPlanShouldDisplayDestroyPlan(t *testing.T) {
	// given
	var stdout, stderr bytes.Buffer

	// when
	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "plan-destroy", awsAccessKey, awsSecretKey)

	//then
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}

func TestOnDestroyShouldDestroyEnvironment(t *testing.T) {
	var stdout, stderr bytes.Buffer

	stdout, stderr = utils.RunCommand(dockerExecPath, "run", "--rm", "-v", mountDir, "-t", imageTag, "destroy", awsAccessKey, awsSecretKey)

	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	log.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
}
