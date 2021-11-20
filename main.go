package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const version = "0.0.0"

func main() {
	args := os.Args

	if len(args) > 1 {
		switch args[1] {
		case "version":
			commandVersion()
		case "exec":
			commandExec(args[2:])
		case "history":
			commandHistory()
		}

	}
}

func commandVersion() {
	fmt.Printf("furoshiki5 version %s\n", version)
}

func commandExec(cmds []string) {
	cmdName := cmds[0]
	args := cmds[1:]

	_, err := exec.LookPath(cmdName)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: init repository

	var out bytes.Buffer
	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	var logs []string
	for {
		line, err := out.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		logs = append(logs, line)
	}

	logHeader := []string{
		fmt.Sprintf("command:     %s", "json"),
		fmt.Sprintf("user:        %s", getUsername()),
		fmt.Sprintf("repoPath:    %s", "repoPath"),
		fmt.Sprintf("projectPath: %s", "projectPath"),
		fmt.Sprintf("gitRevision: %s", "gitRevision"),
		fmt.Sprintf("furoVersion: %s", version),
		fmt.Sprintf("exitCode:    %d", cmd.ProcessState.ExitCode()),
		"---\n",
	}

	fmt.Printf("%s", strings.Join(logHeader, "\n"))
	fmt.Printf("%s", strings.Join(logs, ""))

	// TODO: push repository
}

func getUsername() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return user.Username
}

func commandHistory() {
	fmt.Println("not implemented yet")
}
