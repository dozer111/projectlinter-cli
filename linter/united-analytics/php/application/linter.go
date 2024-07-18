//go:generate go run github.com/dozer111/projectlinter-core/rules/dependency/bump/go_build_generator -c "$PWD"
//go:generate go run github.com/dozer111/projectlinter-core/rules/dependency/substitute/go_build_generator -c "$PWD"
package application

import (
	"embed"
	"fmt"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/common_sets"
	"github.com/dozer111/projectlinter-core/rules"
	"github.com/dozer111/projectlinter-core/rules/dependency"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"os"
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
	rootDirEntries, err := os.ReadDir(l.pathProvider.PathToCaller())
	if err != nil {
		return nil, fmt.Errorf("cannot parse dir %s: %w", l.pathProvider.PathToCaller(), err)
	}

	s_composer := NewComposerSet(l.pathProvider, rootDirEntries)
	sets := []rules.Set{
		s_composer,
		dependency.NewPHPDependencySet(
			l.pathProvider,
			substituteLibraryConfigs,
			bumpLibraryConfigs,
		),
		NewMakefileSet(
			l.pathProvider,
			sources,
			s_composer,
		),
		common_sets.NewEditorConfig(
			l.pathProvider,
			sources,
			"https://your_github.com/projects/your_php_template/browse/.editorconfig",
		),
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

// optionalLeaf - this is commonly used practise for sets in projectlinter
//
// # The main reason is - readability
//
// As for me - it is easier to read the code like
// NewRuleTree
//
//	leaf
//	leaf
//		optionalLeaf
//		leaf
//	leaf
//	optionalLeaf
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
func optionalLeaf(r rules.Rule, children ...rules.RuleTreeLeaf) rules.RuleTreeLeaf {
	return rules.NewOptionalLeaf(r, children...)
}

// leafWithConditions - this is commonly used practise for sets in projectlinter
//
// # The main reason is - readability
//
// As for me - it is easier to read the code like
// NewRuleTree
//
//	leaf
//	leaf
//		leafWithConditions
//		leaf
//	leaf
//	optionalLeaf
//		leaf
//	leaf
//	leafWithConditions
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
func leafWithConditions(r rules.Rule, conditions ...func() bool) rules.RuleTreeLeaf {
	return rules.NewLeafWithConditions(r, conditions)
}
