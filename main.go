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
	"strconv"
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

	gitRevision, err := gitOutput("rev-parse", "HEAD")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: init repository

	logs, exitCode, err := executeCommand(command[0], command[1:]...)
	if err != nil {
		log.Fatal(err)
	}

	logHeader := []string{
		fmt.Sprintf("command:     %s", dumpCommand(command)),
		fmt.Sprintf("user:        %s", getUsername()),
		fmt.Sprintf("repoPath:    %s", "repoPath"),
		fmt.Sprintf("projectPath: %s", "projectPath"),
		fmt.Sprintf("gitRevision: %s", gitRevision),
		fmt.Sprintf("furoVersion: %s", version),
		fmt.Sprintf("exitCode:    %d", exitCode),
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

func gitOutput(command ...string) (string, error) {
	if furo_debug, _ := strconv.ParseBool(os.Getenv("FURO_DEBUG")); furo_debug {
		fmt.Fprintf(os.Stderr, ">>> RUN %s\n", append([]string{"git"}, command...))
	}
	lines, _, err := executeCommand("git", command...)
	if err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func executeCommand(command string, args ...string) ([]string, int, error) {
	var out bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &out

	err := cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		return nil, exitCode, err
	}

	var lines []string
	for {
		line, err := out.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, exitCode, err
		}
		lines = append(lines, line)
	}

	return lines, exitCode, nil
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
