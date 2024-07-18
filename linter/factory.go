package linter

import (
	"errors"
	"fmt"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/configuration"
	unitedAnalyticsGOApplication "github.com/dozer111/projectlinter-cli/app-layout/linter/united-analytics/go/application"
	unitedAnalyticsGOLibrary "github.com/dozer111/projectlinter-cli/app-layout/linter/united-analytics/go/library"
	unitedAnalyticsPHPApplication "github.com/dozer111/projectlinter-cli/app-layout/linter/united-analytics/php/application"
	unitedAnalyticsPHPLibrary "github.com/dozer111/projectlinter-cli/app-layout/linter/united-analytics/php/library"
	unitedCoreGOApplication "github.com/dozer111/projectlinter-cli/app-layout/linter/united-core/go/application"
	unitedCoreGOLibrary "github.com/dozer111/projectlinter-cli/app-layout/linter/united-core/go/library"
	unitedCorePHPApplication "github.com/dozer111/projectlinter-cli/app-layout/linter/united-core/php/application"
	unitedCorePHPLibrary "github.com/dozer111/projectlinter-cli/app-layout/linter/united-core/php/library"
	"github.com/dozer111/projectlinter-core/rules"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"slices"
)

type Unit string

const (
	UnitedCore      Unit = "united-core"
	UnitedAnalytics Unit = "united-analytics"
)

func (u Unit) Valid() bool {
	return slices.Contains[[]string, string](AvailableUnits, string(u))
}

var AvailableUnits = []string{
	string(UnitedCore),
	string(UnitedAnalytics),
}

type LinterFactory struct {
}

var UnknownUnitErr = errors.New("unknown unit")

func (f *LinterFactory) Create(pathProvider *path_provider.PathProvider, cfg *configuration.Configuration) (rules.Linter, error) {
	unitEnum := Unit(cfg.Unit)

	switch true {
	case !unitEnum.Valid():
		return nil, fmt.Errorf("%w: %s", UnknownUnitErr, cfg.Unit)
	case !cfg.Language.Valid():
		return nil, fmt.Errorf("unknown language %s", cfg.Language)
	case !cfg.Mode.Valid():
		return nil, fmt.Errorf("unknown mode %s", cfg.Mode)
	}

	switch unitEnum {
	case UnitedCore:
		return f.createUnitedCoreLinter(pathProvider, cfg), nil
	case UnitedAnalytics:
		return f.createUnitedAnalyticsLinter(pathProvider, cfg), nil
	default:
		return nil, nil
	}
}

func (f *LinterFactory) createUnitedCoreLinter(pathProvider *path_provider.PathProvider, cfg *configuration.Configuration) rules.Linter {
	if cfg.Language.IsGO() {
		if cfg.Mode.IsApplication() {
			return unitedCoreGOApplication.NewLinter(pathProvider)
		}

		return unitedCoreGOLibrary.NewLinter(pathProvider)
	}

	if cfg.Mode.IsApplication() {
		return unitedCorePHPApplication.NewLinter(pathProvider)
	}

	return unitedCorePHPLibrary.NewLinter(pathProvider)
}

func (f *LinterFactory) createUnitedAnalyticsLinter(pathProvider *path_provider.PathProvider, cfg *configuration.Configuration) rules.Linter {
	if cfg.Language.IsGO() {
		if cfg.Mode.IsApplication() {
			return unitedAnalyticsGOApplication.NewLinter(pathProvider)
		}

		return unitedAnalyticsGOLibrary.NewLinter(pathProvider)
	}

	if cfg.Mode.IsApplication() {
		return unitedAnalyticsPHPApplication.NewLinter(pathProvider)
	}

	return unitedAnalyticsPHPLibrary.NewLinter(pathProvider)
}
