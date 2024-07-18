package application

import (
	"fmt"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/configuration"
	pkgTest "github.com/dozer111/projectlinter-cli/app-layout/pkg/test"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	utilTest "github.com/dozer111/projectlinter-core/util/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceLinterRunAllTheRules(t *testing.T) {
	projectDir := utilTest.PathInProjectLinter("testdata/perfect_project")
	pathProvider := path_provider.NewPathProvider(projectDir)

	config, err := configuration.NewParser(fmt.Sprintf("%s/%s", projectDir, configuration.ConfigFileName)).Parse()
	assert.Nil(t, err)
	assert.Equal(t, config.Unit, pkgTest.UnitName())

	l := NewLinter(pathProvider)
	result, err := l.Sets()
	assert.Nil(t, err)

	expectedRuleCount := map[string]int{
		"composer":      21,
		"dependency":    2,
		".editorconfig": 2,
		"makefile":      5,
	}

	for _, set := range result {
		assert.True(t, len(set.Init()) == 0)

		rules := set.Run().Resolve(nil)

		failRules := utilTest.AllSetRulesArePassed(rules)
		assert.True(t, len(failRules) == 0, "Some rules are failed:\n%v", failRules)
		assert.Equal(t, len(rules), expectedRuleCount[set.ID()])
	}
}
