package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
	Long:  `Manage accounts.`,
}

var accountNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new account",
	Long:  `Create a new account.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		teams, _ := cmd.Flags().GetStringSlice("teams")
		output, _ := cmd.Flags().GetString("output")

		fmt.Println("Creating new accounts...")
		return createAccounts(teams, output)
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	accountCmd.AddCommand(accountNewCmd)
	accountNewCmd.Flags().StringSlice("teams", []string{}, "A list of teams to create accounts for")
	accountNewCmd.Flags().String("output", "", "The path to the output directory for the generated accounts")
	accountNewCmd.MarkFlagRequired("teams")
	accountNewCmd.MarkFlagRequired("output")
}

func createAccounts(teams []string, output string) error {
	for _, team := range teams {
		fmt.Printf("Creating account for team %s...\n", team)

		teamPath := filepath.Join(output, team)
		if err := os.MkdirAll(teamPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory for team %s: %w", team, err)
		}

		localsTfPath := filepath.Join(teamPath, "locals.tf")
		if err := renderTemplate("cmd/templates/locals.tf", localsTfPath, pongo2.Context{"Team": team}); err != nil {
			return fmt.Errorf("failed to render locals.tf template for team %s: %w", team, err)
		}
	}
	return nil
}

func renderTemplate(src, dst string, data pongo2.Context) error {
	tplBytes, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	tpl, err := pongo2.FromString(string(tplBytes))
	if err != nil {
		return err
	}

	output, err := tpl.Execute(data)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, []byte(output), 0644)
}
