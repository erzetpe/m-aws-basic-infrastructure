package utils

import (
	"bytes"
	"log"
	"os/exec"
)

func RunCommand(commandWithParams ...string) (bytes.Buffer, bytes.Buffer) {
	var stdout, stderr bytes.Buffer
	dockerRunInit := &exec.Cmd{
		Path:   commandWithParams[0],
		Args:   commandWithParams,
		Stdout: &stdout,
		Stderr: &stderr,
	}
	if err := dockerRunInit.Run(); err != nil {
		log.Println("Error:", err)
	}

	return stdout, stderr
}
