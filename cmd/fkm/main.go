package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ferhatbostanci/fast-key-manager/pkg/github"
	"github.com/ferhatbostanci/fast-key-manager/pkg/gitlab"
	"github.com/ferhatbostanci/fast-key-manager/pkg/ssh"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const Version = "1.0.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "fkm",
		Short: "Fast Key Manager - A tool to manage SSH keys",
		Long: color.New(color.FgCyan).Sprintf(`Fast Key Manager (fkm) is a CLI tool designed to help users manage SSH keys 
on Linux and other Unix-based systems. You can add SSH keys from GitHub or GitLab accounts, 
list existing SSH keys, and remove selected keys from your system.`),
	}

	var addCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a new SSH key",
		Long:  color.New(color.FgCyan).Sprintf(`Add a new SSH key either manually or from GitHub/GitLab accounts.`),
		Run: func(cmd *cobra.Command, args []string) {
			keyManager, err := ssh.NewKeyManager()
			if err != nil {
				color.Red("Error initializing key manager: %v\n", err)
				os.Exit(1)
			}

			if len(args) == 0 && !cmd.Flags().Changed("github") && !cmd.Flags().Changed("gitlab") && !cmd.Flags().Changed("manual") {
				prompt := promptui.Select{
					Label: "Select source for SSH key",
					Items: []string{"Manual Input", "GitHub", "GitLab"},
					Templates: &promptui.SelectTemplates{
						Label:    "{{ . | cyan }}",
						Active:   "\U0001F449 {{ . | cyan }}",
						Inactive: "  {{ . | white }}",
						Selected: "\U0001F44D {{ . | green }}",
					},
				}

				_, result, err := prompt.Run()
				if err != nil {
					color.Red("Prompt failed: %v\n", err)
					return
				}

				switch result {
				case "Manual Input":
					handleManualKey(keyManager)
				case "GitHub":
					promptUsername("GitHub", keyManager, handleGitHubKeys)
				case "GitLab":
					promptUsername("GitLab", keyManager, handleGitLabKeys)
				}
				return
			}

			github, _ := cmd.Flags().GetString("github")
			gitlab, _ := cmd.Flags().GetString("gitlab")
			manual, _ := cmd.Flags().GetBool("manual")

			if github != "" {
				handleGitHubKeys(github, keyManager)
			} else if gitlab != "" {
				handleGitLabKeys(gitlab, keyManager)
			} else if manual {
				handleManualKey(keyManager)
			}
		},
	}

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List all SSH keys",
		Run: func(cmd *cobra.Command, args []string) {
			keyManager, err := ssh.NewKeyManager()
			if err != nil {
				color.Red("Error initializing key manager: %v\n", err)
				os.Exit(1)
			}

			keys, err := keyManager.ListKeys()
			if err != nil {
				color.Red("Error listing keys: %v\n", err)
				os.Exit(1)
			}

			if len(keys) == 0 {
				color.Yellow("No SSH keys found.")
				return
			}

			fmt.Println(color.New(color.FgCyan).Sprint("\n🔑 Your SSH Keys:\n"))
			for i, key := range keys {
				keyInfo := strings.Split(key, " ")
				if len(keyInfo) >= 3 {
					keyType := color.New(color.FgGreen).Sprint(keyInfo[0])
					keyHash := color.New(color.FgYellow).Sprint(keyInfo[1][:20] + "...")
					keyTitle := color.New(color.FgBlue).Sprint(keyInfo[2])
					fmt.Printf("%d. %s %s (%s)\n", i+1, keyType, keyHash, keyTitle)
				} else if len(keyInfo) >= 2 {
					keyType := color.New(color.FgGreen).Sprint(keyInfo[0])
					keyHash := color.New(color.FgYellow).Sprint(keyInfo[1][:20] + "...")
					fmt.Printf("%d. %s %s\n", i+1, keyType, keyHash)
				} else {
					fmt.Printf("%d. %s\n", i+1, key)
				}
			}
			fmt.Println()
		},
	}

	var removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove an SSH key",
		Run: func(cmd *cobra.Command, args []string) {
			keyManager, err := ssh.NewKeyManager()
			if err != nil {
				color.Red("Error initializing key manager: %v\n", err)
				os.Exit(1)
			}

			keys, err := keyManager.ListKeys()
			if err != nil {
				color.Red("Error listing keys: %v\n", err)
				os.Exit(1)
			}

			if len(keys) == 0 {
				color.Yellow("No SSH keys found.")
				return
			}

			var displayKeys []string
			for i, key := range keys {
				keyInfo := strings.Split(key, " ")
				if len(keyInfo) >= 3 {
					keyType := color.New(color.FgGreen).Sprint(keyInfo[0])
					keyHash := color.New(color.FgYellow).Sprint(keyInfo[1][:20] + "...")
					keyTitle := color.New(color.FgBlue).Sprint(keyInfo[2])
					displayKeys = append(displayKeys, fmt.Sprintf("%d. %s %s (%s)", i+1, keyType, keyHash, keyTitle))
				} else if len(keyInfo) >= 2 {
					keyType := color.New(color.FgGreen).Sprint(keyInfo[0])
					keyHash := color.New(color.FgYellow).Sprint(keyInfo[1][:20] + "...")
					displayKeys = append(displayKeys, fmt.Sprintf("%d. %s %s", i+1, keyType, keyHash))
				} else {
					displayKeys = append(displayKeys, fmt.Sprintf("%d. %s", i+1, key))
				}
			}

			prompt := promptui.Select{
				Label: "Select a key to remove",
				Items: displayKeys,
				Templates: &promptui.SelectTemplates{
					Label:    "{{ . | cyan }}",
					Active:   "\U0001F449 {{ . }}",
					Inactive: "  {{ . }}",
					Selected: "\U0001F44D {{ . | red }}",
				},
			}

			index, _, err := prompt.Run()
			if err != nil {
				color.Red("Prompt failed: %v\n", err)
				return
			}

			fmt.Printf("\n%s %s\n", color.New(color.FgYellow).Sprint("⚠️  Warning:"), displayKeys[index])
			fmt.Printf("%s [y/N]: ", color.New(color.FgCyan).Sprint("Are you sure you want to remove this key?"))

			var input string
			fmt.Scanln(&input)

			if strings.ToLower(input) != "y" {
				color.Yellow("Key removal cancelled")
				return
			}

			if err := keyManager.RemoveKey(keys[index]); err != nil {
				color.Red("Error removing key: %v\n", err)
				return
			}

			color.Green("Key removed successfully")
		},
	}

	addCmd.Flags().String("github", "", "GitHub username to fetch SSH keys from")
	addCmd.Flags().String("gitlab", "", "GitLab username to fetch SSH keys from")
	addCmd.Flags().Bool("manual", false, "Manually add an SSH key")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of fkm",
		Run: func(cmd *cobra.Command, args []string) {
			color.Cyan("Fast Key Manager v%s", Version)
		},
	}

	rootCmd.AddCommand(addCmd, listCmd, removeCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		color.Red("%v\n", err)
		os.Exit(1)
	}
}

func promptUsername(service string, keyManager *ssh.KeyManager, handler func(string, *ssh.KeyManager)) {
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("Enter your %s username", service),
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) < 1 {
				return fmt.Errorf("username cannot be empty")
			}
			return nil
		},
	}

	username, err := prompt.Run()
	if err != nil {
		color.Red("Prompt failed: %v\n", err)
		return
	}

	handler(username, keyManager)
}

func handleGitHubKeys(username string, keyManager *ssh.KeyManager) {
	client := github.NewClient()
	keys, err := client.GetUserKeys(username)
	if err != nil {
		color.Red("Error fetching GitHub keys: %v\n", err)
		return
	}

	if len(keys) == 0 {
		color.Yellow("No SSH keys found for GitHub user %s\n", username)
		return
	}

	for _, key := range keys {
		keyName := fmt.Sprintf("github-%s-key-%d", username, key.ID)
		if err := keyManager.AddKey(key.Key, keyName); err != nil {
			color.Red("Error adding key '%d': %v\n", key.ID, err)
		} else {
			color.Green("Added key: %s\n", keyName)
		}
	}
}

func handleGitLabKeys(username string, keyManager *ssh.KeyManager) {
	client := gitlab.NewClient()
	keys, err := client.GetUserKeys(username)
	if err != nil {
		color.Red("Error fetching GitLab keys: %v\n", err)
		return
	}

	if len(keys) == 0 {
		color.Yellow("No SSH keys found for GitLab user %s\n", username)
		return
	}

	for _, key := range keys {
		keyName := fmt.Sprintf("gitlab-%s-key-%d", username, key.ID)
		if err := keyManager.AddKey(key.Key, keyName); err != nil {
			color.Red("Error adding key '%d': %v\n", key.ID, err)
		} else {
			color.Green("Added key: %s\n", keyName)
		}
	}
}

func handleManualKey(keyManager *ssh.KeyManager) {
	commentPrompt := promptui.Prompt{
		Label: "Enter a name/comment for the key",
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) < 1 {
				return fmt.Errorf("comment cannot be empty")
			}
			return nil
		},
	}

	comment, err := commentPrompt.Run()
	if err != nil {
		color.Red("Prompt failed: %v\n", err)
		return
	}

	keyPrompt := promptui.Prompt{
		Label: "Paste your SSH key",
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) < 1 {
				return fmt.Errorf("SSH key cannot be empty")
			}
			return nil
		},
	}

	key, err := keyPrompt.Run()
	if err != nil {
		color.Red("Prompt failed: %v\n", err)
		return
	}

	if err := keyManager.AddKey(strings.TrimSpace(key), strings.TrimSpace(comment)); err != nil {
		color.Red("Error adding key: %v\n", err)
		return
	}

	color.Green("Key added successfully")
} 