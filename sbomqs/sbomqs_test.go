package sbomqs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MockGetBitbucketRepoSlug() string {
	return "MOCK_BITBUCKET_REPO_SLUG"
}

func MockGetBitbucketRepoSlugEmpty() string {
	return ""
}

func MockGetUTCTime() time.Time {
	// Set a fixed time for testing purposes.
	return time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
}

func TestGenSbomqsArgsWithJsonOutputFormat(t *testing.T) {
	// Set SBOMQS_OUTPUT_FORMAT to json for the duration of the test.
	oldEnv := os.Getenv("SBOMQS_OUTPUT_FORMAT")
	defer func() { os.Setenv("SBOMQS_OUTPUT_FORMAT", oldEnv) }()
	os.Setenv("SBOMQS_OUTPUT_FORMAT", "json")

	args := GenSbomqsArgs()

	expectedArgs := "score --json "
	assert.Equal(t, expectedArgs, args)
}

func TestGenSbomqsArgsWithDetailedOutputFormat(t *testing.T) {
	// Set SBOMQS_OUTPUT_FORMAT to json for the duration of the test.
	oldEnv := os.Getenv("SBOMQS_OUTPUT_FORMAT")
	defer func() { os.Setenv("SBOMQS_OUTPUT_FORMAT", oldEnv) }()
	os.Setenv("SBOMQS_OUTPUT_FORMAT", "detailed")

	args := GenSbomqsArgs()

	expectedArgs := "score --detailed "
	assert.Equal(t, expectedArgs, args)
}

func TestGenSbomqsArgsWithNoOutputFormatSet(t *testing.T) {
	// Set SBOMQS_OUTPUT_FORMAT to json for the duration of the test.
	oldEnv := os.Getenv("SBOMQS_OUTPUT_FORMAT")
	defer func() { os.Setenv("SBOMQS_OUTPUT_FORMAT", oldEnv) }()
	os.Unsetenv("SBOMQS_OUTPUT_FORMAT")

	args := GenSbomqsArgs()

	expectedArgs := "score"
	assert.Equal(t, expectedArgs, args)
}
