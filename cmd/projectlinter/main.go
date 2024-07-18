package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/dozer111/projectlinter-cli/app-layout/cmd/projectlinter/cmd"
	"github.com/spf13/cobra"
)

//go:embed version
var localVersion string

func main() {
	currentVersion := semver.MustParse(localVersion)

	rootCmd := &cobra.Command{
		Use:     "projectlinter",
		Version: currentVersion.String(),
		Short:   "static analyzer for project configuration files",
	}

	rootCmd.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd: true,
	}

	setEnvs()
	rootCmd.AddCommand(cmd.CreateCommands(currentVersion)...)

	rootCmd.Execute()
}

func setEnvs() {
	setenv := func(key, value string) {
		if err := os.Setenv(key, value); err != nil {
			log.Fatalf("Failure while set %s: %v", key, err)
		}
	}

	// need for "self-update" command
	setenv("GOPRIVATE", "your_goprivate")
}
