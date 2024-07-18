package run

import (
	"fmt"
	"github.com/dozer111/projectlinter-cli/app-layout/cmd/projectlinter/cmd/run/printer"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/configuration"
	"github.com/dozer111/projectlinter-cli/app-layout/linter"
	"github.com/dozer111/projectlinter-cli/app-layout/pkg/app"
	"github.com/dozer111/projectlinter-core/rules"
	"github.com/dozer111/projectlinter-core/util/painter"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/zoumo/goset"
)

type App struct {
	pathProvider   *path_provider.PathProvider
	currentVersion *semver.Version
	config         *configuration.Configuration
}

func NewApp(
	pathProvider *path_provider.PathProvider,
	currentVersion *semver.Version,
) *App {
	return &App{
		pathProvider:   pathProvider,
		currentVersion: currentVersion,
	}
}

func (c *App) Run() {
	c.checkAppIsLatest()

	if err := c.run(); err != nil {
		fmt.Println("projectlinter run is failed:", err)
		os.Exit(1)
	}
}

func (c *App) runLinter(l rules.Linter, ignore []string) ([]printer.Data, bool) {
	resultSuccess := true
	ignoreAsSet := goset.NewSetFromStrings(ignore)

	sets, err := l.Sets()
	if err != nil {
		fmt.Printf("cannot init linter: %s", err)
		os.Exit(1)
	}

	printData := make([]printer.Data, 0, len(sets))
	for _, set := range sets {
		if ignoreAsSet.Contains(set.ID()) {
			continue
		}

		if initErrs := set.Init(); len(initErrs) > 0 {
			resultSuccess = false
			printData = append(printData, printer.Data{SetID: set.ID(), SetInitErrors: initErrs})
			continue
		}

		setRules := set.Run().Resolve(ignore)
		for _, rule := range setRules {
			if !resultSuccess {
				break
			}

			resultSuccess = rule.IsPassed()
		}

		printData = append(printData, printer.Data{SetID: set.ID(), Rules: setRules})
	}

	return printData, resultSuccess
}

func (c *App) run() error {
	cfg, err := c.parseConfig()
	if err != nil {
		return fmt.Errorf("cannot parse app config: %s", err)
	}

	currentLinter, err := c.createLinters(cfg)
	if err != nil {
		return fmt.Errorf("error while detect linter: %s", err)
	}

	printData, resultSuccess := c.runLinter(currentLinter, cfg.Ignore)

	var output printer.Printer
	switch outputFlag {
	case OutputDefault:
		output = printer.NewDefault()
	case OutputNames:
		output = printer.NewOnlyNames()
	default:
		output = printer.NewDefault()
	}
	fmt.Println(output.Print(printData))

	if !resultSuccess {
		os.Exit(1)
	}

	return nil
}

func (c *App) checkAppIsLatest() {
	// time that should be enough for the goroutine to run and display a message
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	done := make(chan bool)
	go func() {
		latestVersion := pkgApp.FetchLatestVersion()

		if latestVersion.GreaterThan(c.currentVersion) {
			output := []string{
				strings.Repeat("=", 100),
				fmt.Sprintf("= ⚠ %s", "App is outdated."),
				fmt.Sprintf("= Current version: %s", c.currentVersion),
				fmt.Sprintf("= Latest version: %s", latestVersion),
				"Run \"projectlinter self-update\"",
				strings.Repeat("=", 100),
			}

			paint := painter.NewPainter()
			for _, msg := range output {
				fmt.Println(paint.Yellow(msg))
			}
		}

		done <- true
	}()

	select {
	case <-done:
	case <-ticker.C:
		return
	}
}

func (c *App) parseConfig() (*configuration.Configuration, error) {
	pathToConfig := c.pathProvider.PathInCaller(configuration.ConfigFileName)
	appConfigParser := configuration.NewParser(pathToConfig)
	cfg, err := appConfigParser.Parse()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *App) createLinters(cfg *configuration.Configuration) (rules.Linter, error) {
	linterFactory := linter.LinterFactory{}
	l, err := linterFactory.Create(c.pathProvider, cfg)
	if err != nil {
		return nil, err
	}

	return l, nil
}
