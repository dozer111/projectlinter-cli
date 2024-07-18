package run

import (
	"github.com/Masterminds/semver/v3"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"github.com/spf13/cobra"
)

var (
	projectPathFlag string
	outputFlag      = OutputDefault
)

func NewRunCommand(currentVersion *semver.Version) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run analysis of project",
		Long:  "Sets analysis of [GO|PHP] projects in [service|library] mode",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{}, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			pathProvider := path_provider.NewPathProvider(projectPathFlag)
			app := NewApp(pathProvider, currentVersion)
			app.Run()
		},
	}

	cmd.Flags().StringVarP(&projectPathFlag, "project-path", "p", "", "for development purposes.Example: -p=/var/www/myProject")
	cmd.Flags().VarP(&outputFlag, "output", "o", `-o names`)

	return cmd
}
