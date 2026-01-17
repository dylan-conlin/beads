package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestOutputContextFunction(t *testing.T) {
	tests := []struct {
		name          string
		mcpMode       bool
		stealthMode   bool
		ephemeralMode bool
		expectText    []string
		rejectText    []string
	}{
		{
			name:          "CLI Normal (non-ephemeral)",
			mcpMode:       false,
			stealthMode:   false,
			ephemeralMode: false,
			expectText:    []string{"Beads Workflow Context", "bd sync", "git push"},
			rejectText:    []string{"bd sync --flush-only", "--from-main"},
		},
		{
			name:          "CLI Normal (ephemeral)",
			mcpMode:       false,
			stealthMode:   false,
			ephemeralMode: true,
			expectText:    []string{"Beads Workflow Context", "bd sync --from-main", "ephemeral branch"},
			rejectText:    []string{"bd sync --flush-only", "git push"},
		},
		{
			name:          "CLI Stealth",
			mcpMode:       false,
			stealthMode:   true,
			ephemeralMode: false, // stealth mode overrides ephemeral detection
			expectText:    []string{"Beads Workflow Context", "bd sync --flush-only"},
			rejectText:    []string{"git push", "git pull", "git commit", "git status", "git add"},
		},
		{
			name:          "MCP Normal (non-ephemeral)",
			mcpMode:       true,
			stealthMode:   false,
			ephemeralMode: false,
			expectText:    []string{"Beads Issue Tracker Active", "bd sync", "git push"},
			rejectText:    []string{"bd sync --flush-only", "--from-main"},
		},
		{
			name:          "MCP Normal (ephemeral)",
			mcpMode:       true,
			stealthMode:   false,
			ephemeralMode: true,
			expectText:    []string{"Beads Issue Tracker Active", "bd sync --from-main", "ephemeral branch"},
			rejectText:    []string{"bd sync --flush-only", "git push"},
		},
		{
			name:          "MCP Stealth",
			mcpMode:       true,
			stealthMode:   true,
			ephemeralMode: false, // stealth mode overrides ephemeral detection
			expectText:    []string{"Beads Issue Tracker Active", "bd sync --flush-only"},
			rejectText:    []string{"git push", "git pull", "git commit", "git status", "git add"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer stubIsEphemeralBranch(tt.ephemeralMode)()

			var buf bytes.Buffer
			err := outputPrimeContext(&buf, tt.mcpMode, tt.stealthMode)
			if err != nil {
				t.Fatalf("outputPrimeContext failed: %v", err)
			}

			output := buf.String()

			for _, expected := range tt.expectText {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected text not found: %s", expected)
				}
			}

			for _, rejected := range tt.rejectText {
				if strings.Contains(output, rejected) {
					t.Errorf("Unexpected text found: %s", rejected)
				}
			}
		})
	}
}

func TestIsSpawnedContext(t *testing.T) {
	tests := []struct {
		name           string
		createFile     bool
		expectedResult bool
	}{
		{
			name:           "SPAWN_CONTEXT.md exists",
			createFile:     true,
			expectedResult: true,
		},
		{
			name:           "SPAWN_CONTEXT.md does not exist",
			createFile:     false,
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir := t.TempDir()

			// Change to temp directory
			originalDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current directory: %v", err)
			}
			defer os.Chdir(originalDir)

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("Failed to change directory: %v", err)
			}

			// Create SPAWN_CONTEXT.md if test requires it
			if tt.createFile {
				if err := os.WriteFile("SPAWN_CONTEXT.md", []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to create SPAWN_CONTEXT.md: %v", err)
				}
			}

			// Test isSpawnedContext
			result := isSpawnedContext()
			if result != tt.expectedResult {
				t.Errorf("isSpawnedContext() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

// stubIsEphemeralBranch temporarily replaces isEphemeralBranch
// with a stub returning returnValue.
//
// Returns a function to restore the original isEphemeralBranch.
// Usage:
//
//	defer stubIsEphemeralBranch(true)()
func stubIsEphemeralBranch(isEphem bool) func() {
	original := isEphemeralBranch
	isEphemeralBranch = func() bool {
		return isEphem
	}
	return func() {
		isEphemeralBranch = original
	}
}
