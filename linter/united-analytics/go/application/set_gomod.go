package application

import (
	"github.com/Masterminds/semver/v3"
	"github.com/dozer111/projectlinter-core/rules"
	"github.com/dozer111/projectlinter-core/rules/golang/gomod/config"
	"github.com/dozer111/projectlinter-core/rules/golang/gomod/parser"
	"github.com/dozer111/projectlinter-core/rules/golang/gomod/rule"
	"github.com/dozer111/projectlinter-core/util/path_provider"
)

type gomodSet struct {
	pathProvider *path_provider.PathProvider

	gomodConfig *config.Config
}

var _ rules.Set = (*gomodSet)(nil)

func NewGomodSet(pathProvider *path_provider.PathProvider) *gomodSet {
	return &gomodSet{
		pathProvider: pathProvider,
	}
}

func (s *gomodSet) ID() string {
	return "go.mod"
}

func (s *gomodSet) Init() []error {
	gomodParser := parser.NewParser(s.pathProvider.PathToCaller())
	gomodConfig, err := gomodParser.Parse()
	if err != nil {
		return []error{err}
	}

	s.gomodConfig = gomodConfig
	return nil
}

func (s *gomodSet) Run() *rules.RuleTree {
	// https://your_github.com/projects/your_template/browse/go.mod
	goModIsLatest := rule.NewGOVersionIsGreaterEqualRule(
		*semver.MustParse(s.gomodConfig.Modfile.Go.Version),
		*semver.MustParse("1.22"),
	)

	return rules.NewRuleTree(
		leaf(goModIsLatest),
	)
}
