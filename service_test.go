package sctl

import (
	"os"
	"runtime"

	"testing"
)

func TestSubstituteEnvArgs(t *testing.T) {
	testData := []string{"$ENV_KEY", "VALUE", "$NO_KEY", "$SCHMONEY$"}
	expectedResults := []string{"ENV_VAL", "VALUE", "$NO_KEY", "$SCHMONEY$"}
	os.Setenv("ENV_KEY", expectedResults[0])
	var expectedResult string
	for idx, key := range testData {
		expectedResult = expectedResults[idx]
		result := subsituteEnvArg(key)
		if expectedResult != result {
			t.Error(
				"For key", key,
				"Expected:", expectedResult,
				"Found:", result,
			)
		}
	}
}

func TestServiceFilePath(t *testing.T) {
	var servicePath, expectedResult string
	if runtime.GOOS == "windows" {
		servicePath = "home\\user"
		expectedResult = "home\\user\\test-service.json"
	} else {
		servicePath = "home/user"
		expectedResult = "home/user/test-service.json"
	}
	testService, _ := getTestService()
	result := testService.ServiceFilePath(servicePath)
	if result != expectedResult {
		t.Error(
			"Incorrect filepath",
			"Expected:", expectedResult,
			"Found:", result,
		)
	}
}

func TestParseEnvVars(t *testing.T) {
	testService, _ := getTestService()
	expectedEnvVars := []string{
		"ENV_VAR=ENV_VAL", "-e", "ENV_VAR=VALUE", "ENV_VAR=$NO_KEY",
		"ENV_VAR=$SCHMONEY$", "ENV_VAR2", "ENV_VAL",
	}
	var expectedResult string
	for idx, arg := range testService.EnvVars {
		expectedResult = expectedEnvVars[idx]
		result := parseEnvArgs(arg)
		if expectedResult != result {
			t.Error(
				"For key", arg,
				"Expected:", expectedResult,
				"Found:", result,
			)
		}
	}
}

func TestStopCommand(t *testing.T) {
	testService, _ := getTestService()
	expectedResult := "docker service rm test-service"
	command := testService.StopCommand()
	result := command.ToString()
	if expectedResult != result {
		t.Error(
			"For command", command,
			"Expected:", expectedResult,
			"Found:", result,
		)
	}
}

func TestPushCommand(t *testing.T) {
	testService, _ := getTestService()
	expectedResult := "docker push test/image"
	command := testService.PushCommand()
	result := command.ToString()
	if expectedResult != result {
		t.Error(
			"For command", command,
			"Expected:", expectedResult,
			"Found:", result,
		)
	}
}

func TestPullCommand(t *testing.T) {
	testService, _ := getTestService()
	expectedResult := "docker pull test/image"
	command := testService.PullCommand()
	result := command.ToString()
	if expectedResult != result {
		t.Error(
			"For command", command,
			"Expected:", expectedResult,
			"Found:", result,
		)
	}
}

func TestGetKeywordArgs(t *testing.T) {
	testService, expectedService := getTestService()
	var expectedResult string
	for idx, result := range testService.GetKeywordArgs() {
		expectedResult = expectedService.KeywordArgs[idx]
		if expectedResult != result {
			t.Error(
				"Expected:", expectedResult,
				"Found:", result,
			)
		}
	}
}

func TestGetEnvArgs(t *testing.T) {
	testService, expectedService := getTestService()
	var expectedResult string
	results := testService.GetEnvVars()
	for idx, result := range results {
		expectedResult = expectedService.EnvVars[idx]
		if expectedResult != result {
			t.Error(
				"Expected:", expectedResult,
				"Found:", result,
			)
		}
	}
}

func TestRunCommand(t *testing.T) {
	testService, _ := getTestService()
	network := "test-net"
	runCommand := testService.RunCommand(network)
	expectedMain := "docker"
	if expectedMain != runCommand.Main {
		t.Error(
			"Expected:", expectedMain,
			"Found:", runCommand.Main,
		)
	}
	expectedArgs := []string{
		"service", "create", "--name", "test-service", "--network", "test-net",
		"-e", "ENV_VAR=ENV_VAL", "-e", "ENV_VAR=VALUE", "-e", "ENV_VAR=$NO_KEY",
		"-e", "ENV_VAR=$SCHMONEY$", "-e", "ENV_VAR2", "-e", "ENV_VAL",
		"-p", "80:80", "-d", "test/image",
	}
	var expectedArg string
	for idx, arg := range runCommand.Args {
		expectedArg = expectedArgs[idx]
		if expectedArg != arg {
			t.Error(
				"Expected:", expectedArg,
				"Found:", arg,
			)
		}
	}
}

func TestGetServiceDef(t *testing.T) {
	_, expectedService := getTestService()
	testService := Service{
		Name: "test-service",
	}
	err := testService.GetServiceDef("test-data")
	if err != nil {
		t.Error(err.Error())
	}
	foundRunCommand := testService.RunCommand("test").ToString()
	expectedRunCommand := expectedService.RunCommand("test").ToString()
	if foundRunCommand != expectedRunCommand {
		t.Error(
			"Expected:", expectedRunCommand,
			"Found:", foundRunCommand,
		)
	}
}

func getTestService() (Service, Service) {
	os.Setenv("ENV_KEY", "ENV_VAL")
	testService := Service{
		Name:        "test-service",
		Image:       "test/image",
		KeywordArgs: []string{"-p 80:80", "-d"},
		EnvVars: []string{
			"ENV_VAR=$ENV_KEY", "-e", "ENV_VAR=VALUE", "ENV_VAR=$NO_KEY",
			"ENV_VAR=$SCHMONEY$", "ENV_VAR2", "$ENV_KEY",
		},
	}
	expectedService := Service{
		Name:        "test-service",
		Image:       "test/image",
		KeywordArgs: []string{"-p", "80:80", "-d"},
		EnvVars: []string{
			"-e", "ENV_VAR=ENV_VAL", "-e", "ENV_VAR=VALUE", "-e", "ENV_VAR=$NO_KEY",
			"-e", "ENV_VAR=$SCHMONEY$", "-e", "ENV_VAR2", "-e", "ENV_VAL",
		},
	}
	return testService, expectedService
}
