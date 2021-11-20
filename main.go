package main

import (
	"bytes"
	"encoding/json"
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
		default:
			commandHelp()
		}
	} else {
		commandHelp()
	}
}

func commandVersion() {
	fmt.Printf("furoshiki5 version %s\n", version)
}

func commandExec(command []string) {
	if len(command) == 0 {
		commandHelp()
	}

	_, err := exec.LookPath(command[0])
	if err != nil {
		log.Fatal(err)
	}

	// TODO: init repository

	var out bytes.Buffer
	cmd := exec.Command(command[0], command[1:]...)
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
		fmt.Sprintf("command:     %s", dumpCommand(command)),
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

func dumpCommand(command []string) string {
	b, err := json.Marshal(command)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func commandHistory() {
	fmt.Println("not implemented yet")
}

func commandHelp() {
	fmt.Fprintln(os.Stderr, `furo5 exec COMMAND [ARGS...]
furo5 history [pull | show COMMIT | fix]
furo5 version`)
	os.Exit(127)
}
