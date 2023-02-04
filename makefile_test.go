package main_test

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMakefile tests the functionality of the Makefile.
func TestMakefile(t *testing.T) {
	// Not suitable for running parallel due to environment variable modification.

	type testCase struct {
		setup         func(*testing.T)
		withEnv       map[string]string
		executable    string
		parameters    []string
		outputChecker func(*testing.T, string)
		cleanup       func(*testing.T)
	}

	testCases := map[string]testCase{
		"version returns the current git describe tag if no VERSION environment variable is set": testCase{
			executable: "make",
			parameters: []string{"version"},
			outputChecker: func(t *testing.T, got string) {
				currentGitTag, err := exec.Command("git", "describe", "--tags", "--always").Output()
				assert.NoError(t, err)
				assert.Equal(t, string(currentGitTag), got, "make version")
			},
		},
		"version returns the VERSION environment variable when set": testCase{
			withEnv: map[string]string{
				"VERSION": "some-version",
			},
			executable: "make",
			parameters: []string{"version"},
			outputChecker: func(t *testing.T, got string) {
				assert.Equal(t, "some-version\n", got, "make version")
			},
		},
		"version returns git tag with v prefix stripped when in git tag": testCase{
			setup: func(t *testing.T) {
				out, err := exec.Command("git", "tag", "v9999.9.9").Output()
				assert.NoError(t, err, "Creating temp git tag - error")
				assert.Empty(t, out, "Creating temp git tag - output")
			},
			executable: "make",
			parameters: []string{"version"},
			outputChecker: func(t *testing.T, got string) {
				assert.Equal(t, "9999.9.9\n", got, "make version")
			},
			cleanup: func(t *testing.T) {
				out, err := exec.Command("git", "tag", "-d", "v9999.9.9").Output()
				assert.NoError(t, err, "Clean up temp git tag - error")
				assert.Contains(t, string(out), "Deleted tag 'v9999.9.9' (was", "Clean up temp git tag - output")
			},
		},
		"version returns environment variable VERSION with v prefix is not stripped from VERSION environment variable when set": testCase{
			withEnv: map[string]string{
				"VERSION": "v3.2.1",
			},
			executable: "make",
			parameters: []string{"version"},
			outputChecker: func(t *testing.T, got string) {
				assert.Equal(t, "v3.2.1\n", got, "make version")
			},
		},
	}

	for tName, tCase := range testCases {
		tn, tc := tName, tCase

		t.Run(tn, func(t *testing.T) {
			// Not suitable for running parallel due to environment variable modification.
			if tc.setup != nil {
				tc.setup(t)
			}

			for k, v := range tc.withEnv {
				t.Setenv(k, v)
			}

			got, err := exec.Command(tc.executable, tc.parameters...).Output()
			assert.NoError(t, err)
			tc.outputChecker(t, string(got))

			if tc.cleanup != nil {
				tc.cleanup(t)
			}
		})
	}
}
