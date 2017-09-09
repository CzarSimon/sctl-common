package sctl

import (
	"bytes"
	"os/exec"
	"strings"
)

// Command Holds the main command and argument to be executed against os/exec
type Command struct {
	Main string   `json:"main"`
	Args []string `json:"args"`
}

// Execute Executes a given command against os/exec and return the output and potential errors
func (command Command) Execute() (string, error) {
	cmd := exec.Command(command.Main, command.Args...)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return errOut.String(), err
	}
	return out.String(), nil
}

// ToString returns the string representation of the command
func (command Command) ToString() string {
	return command.Main + " " + strings.Join(command.Args, " ")
}

// DockerCommand creates a general docker command
func DockerCommand(args []string) Command {
	return Command{
		Main: "docker",
		Args: args,
	}
}

// MinionCommand Holds a command and the minon it should be sent to
type MinionCommand struct {
	Minion  Node    `json:"node"`
	Command Command `json:"command"`
}
