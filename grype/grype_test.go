package grype

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sbom-utilities/utils"
	"testing"
)

func TestGenGrypeOutputFilename_WithCmdArgs(t *testing.T) {
	orgGrypeArgs := "grype --some other args --file filename.txt"
	expectedGrypeOutputFile := "filename.txt"
	expectedResult := orgGrypeArgs

	grypeArgs, grypeOutputFile := GenGrypeOutputFilename(orgGrypeArgs)

	assert.Equal(t, expectedResult, grypeArgs, "Generated grypeArgs does not match expected result")
	assert.Equal(t, expectedGrypeOutputFile, grypeOutputFile, "Generated grypeOutputFile does not match expected result")
}

func TestGenGrypeOutputFilename_WithEnvVariable(t *testing.T) {
	orgGrypeArgs := "grype --some other args"
	expectedGrypeOutputFile := "env_file.txt"
	expectedGrypeArgs := orgGrypeArgs + " --file " + expectedGrypeOutputFile

	_ = utils.SetEnvVariable("GRYPE_OUTPUT_FILENAME", expectedGrypeOutputFile)
	defer os.Unsetenv("GRYPE_OUTPUT_FILENAME")

	grypeArgs, grypeOutputFile := GenGrypeOutputFilename(orgGrypeArgs)

	assert.Equal(t, expectedGrypeArgs, grypeArgs, "Generated grypeArgs does not match expected result")
	assert.Equal(t, expectedGrypeOutputFile, grypeOutputFile, "Generated grypeOutputFile does not match expected result")
}

func TestGenGrypeOutputFilename_WithoutCmdArgsOrEnvVariable(t *testing.T) {
	os.Setenv("BITBUCKET_REPO_SLUG", "some-bitbucket-slug")
	orgGrypeArgs := "grype --some other args"
	expectedGrypeOutputFilePattern := `grype-scan_[\w-]+_\d{8}-\d{2}-\d{2}-\d{2}\.txt`

	grypeArgs, grypeOutputFile := GenGrypeOutputFilename(orgGrypeArgs)

	assert.Equal(t, orgGrypeArgs+" --file "+grypeOutputFile, grypeArgs, "Generated grypeArgs does not match expected result")
	assert.Regexp(t, expectedGrypeOutputFilePattern, grypeOutputFile, "Generated grypeOutputFile does not match expected pattern")
}
