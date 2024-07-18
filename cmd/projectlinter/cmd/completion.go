package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func NewCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

  $ source <yourprogram completion bash>

  # To load completions for each session, execute once:
  # Linux:
  $ projectlinter completion bash > /etc/bash_completion.d/projectlinter
  # macOS:
  $ projectlinter completion bash > /usr/local/etc/bash_completion.d/projectlinter

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ projectlinter completion zsh > "${fpath[1]}/_projectlinter"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ projectlinter completion fish | source

  # To load completions for each session, execute once:
  $ projectlinter completion fish > ~/.config/fish/completions/projectlinter.fish

PowerShell:

  PS> projectlinter completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> projectlinter completion powershell > projectlinter.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletion(os.Stdout)
			}

			return nil
		},
	}
}
