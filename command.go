package sctl

import (
	"log"
	"os/exec"
	"runtime/debug"
	"strings"
)

// Command Holds the main command and argument to be executed against os/exec
type Command struct {
	Main string   `json:"main"`
	Args []string `json:"args"`
}

// Execute Executes a given command against os/exec and return the output and potential errors
func (command Command) Execute() (string, error) {
	out, err := exec.Command(command.Main, command.Args...).Output()
	if err != nil {
		log.Println("Failed: " + command.ToString())
		log.Println(string(debug.Stack()))
	}
	return string(out), err
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
