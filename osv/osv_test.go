package osv

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

func TestGenOSVArgsWithJsonOutputFormat(t *testing.T) {
	// Set OSV_OUTPUT_FORMAT to json for the duration of the test.
	oldEnv := os.Getenv("OSV_OUTPUT_FORMAT")
	defer func() { os.Setenv("OSV_OUTPUT_FORMAT", oldEnv) }()
	os.Setenv("OSV_OUTPUT_FORMAT", "json")

	args := GenOsvArgs()

	expectedArgs := "--format json "
	assert.Contains(t, args, expectedArgs)
}
