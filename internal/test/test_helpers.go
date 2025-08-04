package test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"golang.org/x/mod/semver"
)

// Helper to read testdata files into string.
func TestAccTestdata(t *testing.T, path string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get current file")
	}
	example, err := os.ReadFile(filepath.Join(filepath.Dir(currentFile), "..", "testdata", path))
	if err != nil {
		t.Fatal(err)
	}
	return string(example)
}

// Helper to read examples files into string.
// path is defined relative to examples directory.
func TestAccExample(t *testing.T, path ...string) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not get current file")
	}
	pathFragments := append([]string{filepath.Dir(currentFile), "..", "..", "examples"}, path...)
	example, err := os.ReadFile(filepath.Join(pathFragments...))
	if err != nil {
		t.Fatal(err)
	}
	return string(example)
}

// Check if a string contains all expected values.
func TestCheckResourceAttrContainsStringsFunc(expected ...string) func(value string) error {
	return func(value string) error {
		for _, e := range expected {
			if !strings.Contains(value, e) {
				return fmt.Errorf("expected manifest to contain %q", e)
			}
		}
		return nil
	}
}

// Check if license is setup in env to enable some tests behind license.
func CheckEnterpriseEnabled(t *testing.T) {
	value, exists := os.LookupEnv("CDK_LICENSE")
	if !exists || value == "" {
		t.Skip("Skipping tests in free mode as it requires a license set on CDK_LICENSE env var")
	}
}

// Check if version meets minimum requirement.
// NOTE: Version has to start with v, e.g. "v1.2.3".
func CheckMinimumVersionRequirement(t *testing.T, version string, minimumVersion string) {
	// If version is not valid we assume it's an unreleased version. (e.g. main-xxxxxx)
	// In this case we assume it's a newer API, so we don't need to skip the tests.
	if semver.IsValid(version) && semver.Compare(version, minimumVersion) < 0 {
		t.Skip("Skipping tests as API doesn't meet mininum version requirements")
	}
}

// Provider configuration pre-checks.
func TestAccPreCheck(t *testing.T) {
	// check that the environment variables are set
	if os.Getenv("CDK_BASE_URL") == "" {
		t.Fatal("CDK_BASE_URL must be set for acceptance tests")
	}
	if os.Getenv("CDK_ADMIN_EMAIL") == "" {
		t.Fatal("CDK_ADMIN_EMAIL must be set for acceptance tests")
	}
	if os.Getenv("CDK_ADMIN_PASSWORD") == "" {
		t.Fatal("CDK_ADMIN_PASSWORD must be set for acceptance tests")
	}
}
