package utils

import (
	//"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunBashCommand(t *testing.T) {
	//given
	command := "echo 'Hello, World!'"

	//when
	output, err := RunBashCommand(command)

	//then
	assert.NoError(t, err, "expected no error but got one")
	expectedOutput := "Hello, World!"
	assert.Equal(t, expectedOutput, output)
}

func TestRunBashCommandFailureCase(t *testing.T) {
	//given
	command := "'Hello, World!'"

	//when
	output, err := RunBashCommand(command)

	//then
	assert.Error(t, err, "expexted an error but did not get one")
	expectedOutput := ""
	assert.Equal(t, expectedOutput, output)
}

func TestCheckEnvVarThatExists(t *testing.T) {
	//given
	os.Setenv("TEST_ENV_VAR", "1")
	
	//when
	result := CheckEnvVar("TEST_ENV_VAR")

	//then
	assert.True(t, result)
}

func TestCheckEnvVarThatDoesNotExists(t *testing.T) {
	//when
	result := CheckEnvVar("THIS_ENV_VAR_DOES_NOT_EXIST")

	//then
	assert.False(t, result)
}

func TestEnvVarIsTrue(t *testing.T) {
	//given
	os.Setenv("DO_SOMETHING", "true")

	//when & then
	assert.True(t, CheckIfEnvVarIsTrue("DO_SOMETHING"))
}

func TestEnvVarIsNotTrue(t *testing.T) {
	//given
	os.Setenv("DO_SOMETHING", "false")

	//when & then
	assert.False(t, CheckIfEnvVarIsTrue("DO_SOMETHING"))
}

func TestCheckIfFileExists(t *testing.T) {
	//given
	tempDirName, err := os.MkdirTemp(".", "temp")
	defer os.Remove(tempDirName)
	assert.NoError(t, err, "an error occured when trying to create a temp dir")

	assert.True(t, CheckFileExists(tempDirName))
}

func TestCheckIfFileDoesNotExists(t *testing.T) {
	assert.False(t, CheckFileExists("thisDirDoesNotExist"))
}
