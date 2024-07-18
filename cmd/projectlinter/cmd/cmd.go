package cmd

import (
	"github.com/Masterminds/semver/v3"
	"github.com/dozer111/projectlinter-cli/app-layout/cmd/projectlinter/cmd/run"
	"github.com/spf13/cobra"
)

func CreateCommands(currentVersion *semver.Version) []*cobra.Command {
	return []*cobra.Command{
		NewCompletionCmd(),
		run.NewRunCommand(currentVersion),
		NewSelfUpdateCommand(currentVersion),
		NewInitCommand(),
	}
}
