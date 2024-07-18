package cmd

import (
	"fmt"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/configuration"
	"github.com/dozer111/projectlinter-cli/app-layout/linter"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"os"
	"strings"

	"github.com/1set/gut/yos"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Example: `projectlinter init`,
		Short:   "configure project to work with linter",
		Run: func(cmd *cobra.Command, args []string) {
			pathProvider := path_provider.NewPathProvider("")
			configurationFileName := configuration.ConfigFileName
			if yos.ExistFile(configurationFileName) {
				fmt.Printf("%s already exists\n", configurationFileName)
				return
			}

			var language configuration.Language

			switch true {
			case yos.ExistFile(pathProvider.PathInCaller("composer.json")):
				language = configuration.LanguagePHP
			case yos.ExistFile(pathProvider.PathInCaller("go.mod")):
				language = configuration.LanguageGO
			}

			if language != "" {
				fmt.Printf("The language is automatically identified: %s\n", language)
			} else {
				prompt := promptui.Select{
					Label: "Select programming language",
					Items: configuration.AvailableLanguages,
				}

				_, result, err := prompt.Run()
				if err != nil {
					os.Exit(1)
				}

				language = configuration.Language(result)
			}

			prompt := promptui.Select{
				Label: "Select mode",
				Items: configuration.AvailableMods,
			}

			_, result, err := prompt.Run()
			if err != nil {
				os.Exit(1)
			}

			mode := configuration.Mode(result)

			var availableUnits []string

			availableUnits = make([]string, 0, len(linter.AvailableUnits))
			for _, u := range linter.AvailableUnits {
				availableUnits = append(availableUnits, u)
			}

			prompt = promptui.Select{
				Label: "Select unit",
				Items: availableUnits,
			}

			_, result, err = prompt.Run()
			if err != nil {
				os.Exit(1)
			}

			unit := strings.TrimPrefix(result, "(c) ")

			config := configuration.Configuration{
				Unit:     unit,
				Language: language,
				Mode:     mode,
			}

			file, err := os.Create(pathProvider.PathInCaller(configurationFileName))
			if err != nil {
				fmt.Println("Error creating file:", err)
				os.Exit(1)
			}
			defer file.Close()

			b, err := yaml.Marshal(config)
			if err != nil {
				fmt.Println("Cannot marshal config to yaml")
				os.Exit(1)
			}

			_, err = file.Write(b)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				os.Exit(1)
			}

			fmt.Printf("Generate %s\n", configurationFileName)
		},
	}

	return cmd
}
