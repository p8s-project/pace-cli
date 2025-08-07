package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Vezia/vez-cli/internal/llm"
	"github.com/Vezia/vez-cli/internal/loader"
	"github.com/Vezia/vez-cli/internal/types"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for a catalog",
	Long:  `Generate documentation for a catalog.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		catalogPath, _ := cmd.Flags().GetString("catalog")
		output, _ := cmd.Flags().GetString("output")

		fmt.Println("Generating documentation...")
		return generateDocs(catalogPath, output)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.Flags().String("catalog", "", "Path to the catalog directory (required)")
	docsCmd.Flags().String("output", "", "Path to the output directory for the generated documentation (required)")
	docsCmd.MarkFlagRequired("catalog")
	docsCmd.MarkFlagRequired("output")
}

func generateDocs(catalogPath, output string) error {
	files, err := os.ReadDir(catalogPath)
	if err != nil {
		return fmt.Errorf("failed to read catalog directory: %w", err)
	}

	catalog := &types.Catalog{
		Resources: make(map[string]types.ResourceSpec),
	}

	for _, file := range files {
		if !file.IsDir() {
			path := filepath.Join(catalogPath, file.Name())
			c, err := loader.LoadCatalog(path)
			if err != nil {
				return err
			}
			for name, resource := range c.Resources {
				catalog.Resources[name] = resource
			}
		}
	}

	llmClient := llm.New("http://localhost:11434")

	for name, resource := range catalog.Resources {
		fmt.Printf("Generating documentation for %s...\n", name)

		absCatalogPath, err := filepath.Abs(catalogPath)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for catalog: %w", err)
		}
		mainTfPath := filepath.Join(absCatalogPath, "..", resource.Source, "main.tf")
		mainTf, err := os.ReadFile(mainTfPath)
		if err != nil {
			return fmt.Errorf("failed to read main.tf for module %s: %w", name, err)
		}

		prompt := fmt.Sprintf("Please provide a one-paragraph description of the following Terraform module:\n\n%s", string(mainTf))
		req := &llm.GenerateRequest{
			Model:  "tinyllama",
			Prompt: prompt,
		}
		resp, err := llmClient.Generate(req)
		if err != nil {
			return fmt.Errorf("failed to generate description for module %s: %w", name, err)
		}

		outputPath := filepath.Join(output, fmt.Sprintf("%s.md", name))
		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		f, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()

		if _, err := f.WriteString(fmt.Sprintf("# %s\n\n", name)); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
		if _, err := f.WriteString(fmt.Sprintf("%s\n\n", resp.Response)); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}

		cmd := exec.Command("terraform-docs", "markdown", filepath.Join(filepath.Dir(catalogPath), resource.Source))
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run terraform-docs: %w", err)
		}

		if _, err := f.Write(out.Bytes()); err != nil {
			return fmt.Errorf("failed to write to output file: %w", err)
		}
	}

	return nil
}
