package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Vezia/vez-cli/internal/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Manage catalogs",
	Long:  `Manage catalogs.`,
}

var catalogGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a catalog from a directory of Terraform modules",
	Long:  `Generate a catalog from a directory of Terraform modules.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		from, _ := cmd.Flags().GetString("from")
		output, _ := cmd.Flags().GetString("output")
		sourceType, _ := cmd.Flags().GetString("source-type")
		gitBaseURL, _ := cmd.Flags().GetString("git-base-url")

		if sourceType == "git" && gitBaseURL == "" {
			return fmt.Errorf("--git-base-url is required when --source-type is git")
		}

		return generateCatalog(from, output, sourceType, gitBaseURL)
	},
}

func init() {
	rootCmd.AddCommand(catalogCmd)
	catalogCmd.AddCommand(catalogGenerateCmd)
	catalogGenerateCmd.Flags().String("from", "", "The path to the directory of Terraform modules")
	catalogGenerateCmd.Flags().String("output", "", "The path to the output directory for the generated catalog")
	catalogGenerateCmd.Flags().String("source-type", "local", "The source type of the modules (local or git)")
	catalogGenerateCmd.Flags().String("git-base-url", "", "The base URL of the Git repository")
	catalogGenerateCmd.MarkFlagRequired("from")
	catalogGenerateCmd.MarkFlagRequired("output")
}

func generateCatalog(from, output, sourceType, gitBaseURL string) error {
	return filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			files, err := filepath.Glob(filepath.Join(path, "*.tf"))
			if err != nil {
				return err
			}
			if len(files) > 0 {
				return generateCatalogFile(path, output, sourceType, gitBaseURL)
			}
		}
		return nil
	})
}

func generateCatalogFile(modulePath, outputDir, sourceType, gitBaseURL string) error {
	cmd := exec.Command("terraform-docs", "json", modulePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run terraform-docs: %w", err)
	}

	var moduleDoc struct {
		Inputs  []map[string]interface{} `json:"inputs"`
		Outputs []map[string]interface{} `json:"outputs"`
	}
	if err := json.Unmarshal(out.Bytes(), &moduleDoc); err != nil {
		return fmt.Errorf("failed to parse terraform-docs output: %w", err)
	}

	var inputs []types.InputSpec
	for _, input := range moduleDoc.Inputs {
		inputs = append(inputs, types.InputSpec{
			From:     input["name"].(string),
			To:       input["name"].(string),
			Required: input["required"].(bool),
			Default:  input["default"],
		})
	}

	var outputs []types.OutputSpec
	for _, output := range moduleDoc.Outputs {
		outputs = append(outputs, types.OutputSpec{
			From: output["name"].(string),
			To:   output["name"].(string),
		})
	}

	moduleName := filepath.Base(modulePath)
	var source string
	if sourceType == "local" {
		source = fmt.Sprintf("./%s", modulePath)
	} else {
		source = fmt.Sprintf("git::%s//%s", gitBaseURL, moduleName)
	}

	catalog := types.Catalog{
		Resources: map[string]types.ResourceSpec{
			moduleName: {
				Source:  source,
				Version: "0.0.0",
				Inputs:  inputs,
				Outputs: outputs,
			},
		},
	}

	data, err := yaml.Marshal(catalog)
	if err != nil {
		return fmt.Errorf("failed to marshal catalog: %w", err)
	}

	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.yaml", moduleName))
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write catalog file: %w", err)
	}

	fmt.Printf("Successfully generated catalog file %s\n", outputPath)
	return nil
}
