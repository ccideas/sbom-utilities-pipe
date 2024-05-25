package osv

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sbom-utilities/utils"
	"testing"
	"sbom-utilities/osv"
)

func TestGenOSVArgsWithDefaultArgs(t *testing.T) {
	// Set OSV_OUTPUT_FORMAT to json for the duration of the test.
	oldEnv := os.Getenv("OSV_OUTPUT_FORMAT")
	defer func() { os.Setenv("OSV_OUTPUT_FORMAT", oldEnv) }()
	os.Setenv("OSV_OUTPUT_FORMAT", "json")

	args := osv.GenOsvArgs()

	expectedArgs := "scan --format table"
	assert.Contains(t, args, expectedArgs)
}

func TestGenOsvOutputFilename_WithCmdArgs(t *testing.T) {
	orgOsvArgs := "osv-scanner --some other args --output filename.txt"
	expectedOsvOutputFile := "filename.txt"
	expectedResult := orgOsvArgs

	osvArgs, osvOutputFile := osv.GenOsvOutputFilename(orgOsvArgs)

	assert.Equal(t, expectedResult, osvArgs, "Generated osvArgs does not match expected result")
	assert.Equal(t, expectedOsvOutputFile, osvOutputFile, "Generated osvOutputFile does not match expected result")
}

func TestGenOsvOutputFilename_WithEnvVariable(t *testing.T) {
	orgOsvArgs := "osv --some other args"
	expectedOsvOutputFile := "env_file.txt"
	expectedOsvArgs := orgOsvArgs + " --output " + expectedOsvOutputFile

	_ = utils.SetEnvVariable("OSV_OUTPUT_FILENAME", expectedOsvOutputFile)
	defer os.Unsetenv("OSV_OUTPUT_FILENAME")

	osvArgs, osvOutputFile := osv.GenOsvOutputFilename(orgOsvArgs)

	assert.Equal(t, expectedOsvArgs, osvArgs, "Generated osvArgs does not match expected result")
	assert.Equal(t, expectedOsvOutputFile, osvOutputFile, "Generated osvOutputFile does not match expected result")
}

func TestGenOsvOutputFilename_WithoutCmdArgsOrEnvVariable(t *testing.T) {
	os.Setenv("BITBUCKET_REPO_SLUG", "some-bitbucket-slug")
	orgOsvArgs := "osv --some other args"
	expectedOsvOutputFilePattern := `osv-scan_[\w-]+_\d{8}-\d{2}-\d{2}-\d{2}\.txt`

	osvArgs, osvOutputFile := osv.GenOsvOutputFilename(orgOsvArgs)

	assert.Equal(t, orgOsvArgs+" --output "+osvOutputFile, osvArgs, "Generated osvArgs does not match expected result")
	assert.Regexp(t, expectedOsvOutputFilePattern, osvOutputFile, "Generated osvOutputFile does not match expected pattern")
}
