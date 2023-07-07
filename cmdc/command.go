package cmdc

import (
	"bytes"
	"os"
	"os/exec"
)

func Shell(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

func ShellBytes(command string) ([]byte, error) {
	cmd := exec.Command("sh", "-c", command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}
func ShellString(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return string(stdout.Bytes()), nil
}

func Bash(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

func BashBytes(command string) ([]byte, error) {
	cmd := exec.Command("bash", "-c", command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return stdout.Bytes(), nil
}
func BashString(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return string(stdout.Bytes()), nil
}
