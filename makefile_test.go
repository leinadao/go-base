package main_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMakefile tests the functionality of the Makefile.
//
// Not suitable for running parallel due to environment variable modification.
// --no-print-directory is used to resolve behavior when run in the pipeline.
func TestMakefile(t *testing.T) {
	const (
		// envNameTestMakefileRunning is the name of an environment variable set when this
		// test is running which will prevent the test from recursively running whilst
		// testing the makefile commands.
		envNameTestMakefileRunning = "TEST_MAKEFILE_RUNNING"
	)

	if _, set := os.LookupEnv(envNameTestMakefileRunning); set != true {
		type testCase struct {
			setup         func(*testing.T)
			withEnv       map[string]string
			executable    string
			parameters    []string
			outputChecker func(*testing.T, string)
			cleanup       func(*testing.T)
		}

		testCases := map[string]testCase{
			"version returns the current git describe tag if no VERSION environment variable is set": {
				executable: "make",
				parameters: []string{"--no-print-directory", "version"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					currentGitTag, err := exec.Command("git", "describe", "--tags", "--always").Output()
					assert.NoError(t, err, "Checking current git tag")
					assert.Equal(t, string(currentGitTag), got, "'make version' output")
				},
			},
			"version returns the VERSION environment variable when set": {
				withEnv: map[string]string{
					"VERSION": "some-version",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "version"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Equal(t, "some-version\n", got, "'make version' output")
				},
			},
			"version returns git tag with v prefix stripped when in git tag": {
				setup: func(t *testing.T) {
					t.Helper()
					out, err := exec.Command("git", "tag", "v9999.9.9").Output()
					assert.NoError(t, err, "Creating temp git tag - error")
					assert.Empty(t, out, "Creating temp git tag - output")
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "version"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Equal(t, "9999.9.9\n", got, "'make version' output")
				},
				cleanup: func(t *testing.T) {
					t.Helper()
					out, err := exec.Command("git", "tag", "-d", "v9999.9.9").Output()
					assert.NoError(t, err, "Clean up temp git tag - error")
					assert.Contains(t, string(out), "Deleted tag 'v9999.9.9' (was", "Clean up temp git tag - output")
				},
			},
			"version returns environment variable VERSION with v prefix is not stripped from VERSION environment variable when set": {
				withEnv: map[string]string{
					"VERSION": "v3.2.1",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "version"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Equal(t, "v3.2.1\n", got, "'make version' output")
				},
			},
			"deps runs a go mod tidy command": {
				executable: "make",
				parameters: []string{"--no-print-directory", "deps"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Equal(t, "go mod tidy\n", got, "'make deps' output")
				},
			},
			"test runs a go test command": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Regexp(t, "go test.", got, "'make test' output")
				},
			},
			"test uses dynamic tag to support Confluent kafka go use on multiple devices": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Regexp(t, "-tags [^-]*dynamic", got, "'make test' output")
				},
			},
			"test uses unit tag to run only unit tests": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Regexp(t, "-tags [^-]*unit", got, "'make test' output")
				},
			},
			"test uses a non-cacheable count flag so that no test caching occurs": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Regexp(t, "-count=[0-9]*", got, "'make test' output")
				},
			},
			"test uses coverage.out value for coverprofile flag": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Contains(t, got, "-coverprofile coverage.out", "'make test' output")
				},
			},
			"test uses count value for covermode flag": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Contains(t, got, "-covermode count", "'make test' output")
				},
			},
			"test uses comma separated go list output for coverpkg flag": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					goListCSV, err := exec.Command("bash", "-c", "go list ./... | tr '\n' ','").Output()
					assert.NoError(t, err, "Checking go list expected output - error")

					assert.Contains(t, got, "-coverpkg="+string(goListCSV), "'make test' output")
				},
			},
			"test uses verbose output flag": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Contains(t, got, "-v", "'make test' output")
				},
			},
			"test runs for all paths and subpaths": {
				withEnv: map[string]string{
					envNameTestMakefileRunning: "true",
				},
				executable: "make",
				parameters: []string{"--no-print-directory", "test"},
				outputChecker: func(t *testing.T, got string) {
					t.Helper()
					assert.Regexp(t, `\Q./...\E`, got, "'make test' output")
				},
			},
		}

		for tName, tCase := range testCases {
			tn, tc := tName, tCase

			t.Run(tn, func(t *testing.T) {
				t.Helper()
				// Not suitable for running parallel due to environment variable modification.
				if tc.setup != nil {
					tc.setup(t)
				}

				for k, v := range tc.withEnv {
					t.Setenv(k, v)
				}

				got, err := exec.Command(tc.executable, tc.parameters...).Output() //nolint:gosec // Only test parameters.
				assert.NoError(t, err)
				tc.outputChecker(t, string(got))

				if tc.cleanup != nil {
					tc.cleanup(t)
				}
			})
		}
	}
}
