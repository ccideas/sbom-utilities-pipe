package version

import (
	"testing"
)

// TestGetModuleVersion verifies that GetModuleVersion returns the correct version.
func TestGetModuleVersion(t *testing.T) {
	expectedVersion := "v1.5.0"
	actualVersion := GetModuleVersion()

	if actualVersion != expectedVersion {
		t.Errorf("GetModuleVersion() = %s; want %s", actualVersion, expectedVersion)
	}
}
