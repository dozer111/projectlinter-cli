package cmd

import (
	"fmt"
	"github.com/dozer111/projectlinter-cli/app-layout/pkg/app"
	"log"
	"os/exec"

	"github.com/Masterminds/semver/v3"
	"github.com/spf13/cobra"
)

func NewSelfUpdateCommand(currentVersion *semver.Version) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "self-update",
		Example: `projectlinter self-update`,
		Short:   "up projectlinter(aka composer self-update)",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			latestVersion := pkgApp.FetchLatestVersion()
			fmt.Printf("Current version is: %s\nLatest version is: %s\n", currentVersion, latestVersion)

			if latestVersion.GreaterThan(currentVersion) {
				fmt.Printf("Upgrading to %s...\n", latestVersion)
				command := exec.Command("go",
					"install",
					// TODO chage it
					"your_git.com/projectlinter/cmd/projectlinter@latest",
				)
				out, err := command.CombinedOutput()
				if err != nil {
					log.Panicf("cannot update projectlinter to %s\nOutput: %s\nError: %v", latestVersion, out, err)
				}
				fmt.Println("Update success.")

				return
			}

			fmt.Println("Current version is latest. Nothing to update")
		},
	}

	return cmd
}
