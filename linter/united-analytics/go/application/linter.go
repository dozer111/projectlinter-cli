//go:generate go run github.com/dozer111/projectlinter-core/rules/dependency/bump/go_build_generator -c "$PWD"
//go:generate go run github.com/dozer111/projectlinter-core/rules/dependency/substitute/go_build_generator -c "$PWD"

package application

import (
	"embed"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/common_sets"
	"github.com/dozer111/projectlinter-core/rules"
	"github.com/dozer111/projectlinter-core/rules/dependency"
	"github.com/dozer111/projectlinter-core/util/path_provider"
)

//go:embed source/*
var sources embed.FS

type Linter struct {
	pathProvider *path_provider.PathProvider
}

var _ rules.Linter = (*Linter)(nil)

func NewLinter(provider *path_provider.PathProvider) *Linter {
	return &Linter{pathProvider: provider}
}

func (l *Linter) Sets() ([]rules.Set, error) {
	sets := []rules.Set{
		NewGomodSet(l.pathProvider),
		dependency.NewGolangDependencySet(
			l.pathProvider,
			substituteLibraryConfigs,
			bumpLibraryConfigs,
		),
		common_sets.NewEditorConfig(
			l.pathProvider,
			sources,
			"https://your_github.com/projects/your_golang_application_template/browse/.editorconfig",
		),
		NewMakefileSet(l.pathProvider, sources),
	}

	return sets, nil
}

// leaf - this is commonly used practise for sets in projectlinter
//
// # The main reason is - readability
//
// As for me - it is easier to read the code like
// NewRuleTree
//
//	leaf
//	leaf
//		leaf
//		leaf
//	leaf
//	leaf
//		leaf
//	leaf
//
// instead of
//
// NewRuleTree
//
//	rules.NewLeaf
//	rules.NewLeaf
//		rules.NewLeaf
//		rules.NewLeaf
//	rules.NewLeaf
//
// ...
func leaf(r rules.Rule, children ...rules.RuleTreeLeaf) rules.RuleTreeLeaf {
	return rules.NewLeaf(r, children...)
}
