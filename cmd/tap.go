package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Vezia/vez-cli/internal/config"
	"github.com/spf13/cobra"
)

var tapCmd = &cobra.Command{
	Use:   "tap",
	Short: "Manage taps",
	Long:  `Manage taps.`,
}

var tapAddCmd = &cobra.Command{
	Use:   "add [name] [url]",
	Short: "Add a tap",
	Long:  `Add a tap.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return addTap(args[0], args[1])
	},
}

var tapRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a tap",
	Long:  `Remove a tap.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeTap(args[0])
	},
}

var tapListCmd = &cobra.Command{
	Use:   "list",
	Short: "List taps",
	Long:  `List taps.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listTaps()
	},
}

func init() {
	rootCmd.AddCommand(tapCmd)
	tapCmd.AddCommand(tapAddCmd)
	tapCmd.AddCommand(tapRemoveCmd)
	tapCmd.AddCommand(tapListCmd)
}

func addTap(name, url string) error {
	cfgPath, err := getConfigPath()
	if err != nil {
		return err
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}

	for _, tap := range cfg.Taps {
		if tap.Name == name {
			return fmt.Errorf("tap '%s' already exists", name)
		}
	}

	cfg.Taps = append(cfg.Taps, config.Tap{Name: name, URL: url})
	if cfg.ActiveTap == "" {
		cfg.ActiveTap = name
	}

	if err := config.Save(cfgPath, cfg); err != nil {
		return err
	}

	fmt.Printf("Successfully added tap '%s'\n", name)
	return nil
}

func removeTap(name string) error {
	cfgPath, err := getConfigPath()
	if err != nil {
		return err
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}

	var newTaps []config.Tap
	var found bool
	for _, tap := range cfg.Taps {
		if tap.Name == name {
			found = true
		} else {
			newTaps = append(newTaps, tap)
		}
	}

	if !found {
		return fmt.Errorf("tap '%s' not found", name)
	}

	cfg.Taps = newTaps
	if cfg.ActiveTap == name {
		if len(cfg.Taps) > 0 {
			cfg.ActiveTap = cfg.Taps[0].Name
		} else {
			cfg.ActiveTap = ""
		}
	}

	if err := config.Save(cfgPath, cfg); err != nil {
		return err
	}

	fmt.Printf("Successfully removed tap '%s'\n", name)
	return nil
}

func listTaps() error {
	cfgPath, err := getConfigPath()
	if err != nil {
		return err
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return err
	}

	for _, tap := range cfg.Taps {
		if tap.Name == cfg.ActiveTap {
			fmt.Printf("* %s (%s)\n", tap.Name, tap.URL)
		} else {
			fmt.Printf("  %s (%s)\n", tap.Name, tap.URL)
		}
	}

	return nil
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".pace", "config.yaml"), nil
}
