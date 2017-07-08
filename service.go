package sctl

import "strings"

// Service Holds service configuration info
type Service struct {
	Name        string   `json:"name"`
	Image       string   `json:"image"`
	EnvArgs     []string `json:"envArgs"`
	KeywordArgs []string `json:"keywordArgs"`
}

// PushCommand creates a command to push a docker image
func (service Service) PushCommand() Command {
	return Command{
		Main: "docker",
		Args: []string{"push", service.Image},
	}
}

// PullCommand creates a command to pull a docker image
func (service Service) PullCommand() Command {
	return Command{
		Main: "docker",
		Args: []string{"pull", service.Image},
	}
}

// StopCommand creates a command to pull a docker image
func (service Service) StopCommand() Command {
	return Command{
		Main: "docker",
		Args: []string{"service", "rm", service.Name},
	}
}

// RunCommand creates a command to start a servcie
func (service Service) RunCommand(network string) Command {
	args := []string{"service", "create", "--name", service.Name, "--network", network}
	args = append(args, service.GetEnvArgs()...)
	args = append(args, service.GetKeywordArgs()...)
	return Command{
		Main: "docker",
		Args: append(args, service.Image),
	}
}

// GetEnvArgs Formats the environment arguments in a service definition to be runnable
func (service Service) GetEnvArgs() []string {
	args := make([]string, 0)
	EnvFlag := "-e"
	for _, arg := range service.EnvArgs {
		if arg != EnvFlag {
			args = append(args, EnvFlag, arg)
		}
	}
	return args
}

// GetKeywordArgs Formats the keyword arguments in a service definition to be runnable
func (service Service) GetKeywordArgs() []string {
	args := make([]string, 0)
	for _, arg := range service.KeywordArgs {
		args = append(args, strings.Fields(arg)...)
	}
	return args
}

// DockerCommand creates a general docker command
func DockerCommand(args []string) Command {
	return Command{
		Main: "docker",
		Args: args,
	}
}
