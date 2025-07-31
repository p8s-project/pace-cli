package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCmd(t *testing.T) {
	testCases, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("failed to read testdata directory: %v", err)
	}

	for _, tc := range testCases {
		if !tc.IsDir() {
			continue
		}
		t.Run(tc.Name(), func(t *testing.T) {
			testDir := filepath.Join("testdata", tc.Name())
			appFile := filepath.Join(testDir, "app.yaml")
			catalogFile := filepath.Join(testDir, "catalog.yaml")
			goldenDir := filepath.Join(testDir, "golden")
			outputDir := t.TempDir()

			// Reset flags for each test run
			rootCmd.ResetFlags()
			generateCmd.ResetFlags()
			initGenerateCmdFlags() // Re-initialize flags

			// Mock the CLI arguments
			args := []string{
				"--app-file", appFile,
				"--catalog", catalogFile,
				"--output-dir", outputDir,
			}
			rootCmd.SetArgs(append([]string{"generate"}, args...))

			// Execute the command
			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("command execution failed: %v", err)
			}

			// Compare the output with the golden files
			goldenFiles, err := os.ReadDir(goldenDir)
			if err != nil {
				t.Fatalf("failed to read golden directory: %v", err)
			}

			for _, gf := range goldenFiles {
				generatedPath := filepath.Join(outputDir, gf.Name())
				goldenPath := filepath.Join(goldenDir, gf.Name())

				generatedBytes, err := os.ReadFile(generatedPath)
				if err != nil {
					t.Fatalf("failed to read generated file %s: %v", generatedPath, err)
				}

				goldenBytes, err := os.ReadFile(goldenPath)
				if err != nil {
					t.Fatalf("failed to read golden file %s: %v", goldenPath, err)
				}

				// Normalize line endings for comparison
				generated := strings.ReplaceAll(string(generatedBytes), "\r\n", "\n")
				golden := strings.ReplaceAll(string(goldenBytes), "\r\n", "\n")

				if generated != golden {
					t.Errorf("generated output does not match golden file %s", gf.Name())
					t.Logf("--- GENERATED ---\n%s", generated)
					t.Logf("---- GOLDEN ----\n%s", golden)
				}
			}
		})
	}
}

// initGenerateCmdFlags is a helper to re-initialize flags for each test run.
func initGenerateCmdFlags() {
	generateCmd.Flags().StringVarP(&appFilePath, "app-file", "a", "", "Path to the app.yaml file (required)")
	generateCmd.Flags().StringVarP(&catalogPath, "catalog", "c", "", "Path to the catalog.yaml file (required)")
	generateCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Directory to output the generated Terraform files (required)")
	generateCmd.MarkFlagRequired("app-file")
	generateCmd.MarkFlagRequired("catalog")
	generateCmd.MarkFlagRequired("output-dir")
}
