package application

import (
	"fmt"
	projectlinterCliConfiguration "github.com/dozer111/projectlinter-cli/app-layout/internal/configuration"
	pkgTest "github.com/dozer111/projectlinter-cli/app-layout/pkg/test"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	utilTest "github.com/dozer111/projectlinter-core/util/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceLinterRunAllTheRules(t *testing.T) {
	projectDir := utilTest.PathInProjectLinter("testdata/perfect_project")
	pathProvider := path_provider.NewPathProvider(projectDir)

	config, err := projectlinterCliConfiguration.NewParser(fmt.Sprintf("%s/%s", projectDir, projectlinterCliConfiguration.ConfigFileName)).Parse()
	require.Nil(t, err)
	assert.Equal(t, config.Unit, pkgTest.UnitName())

	l := NewLinter(pathProvider)
	result, err := l.Sets()
	require.Nil(t, err)

	expectedRuleCount := map[string]int{
		"go.mod":        1,
		"dependency":    2,
		".editorconfig": 2,
		"makefile":      2,
	}

	for _, set := range result {
		require.True(t, len(set.Init()) == 0)

		rules := set.Run().Resolve(nil)

		failRules := utilTest.AllSetRulesArePassed(rules)
		assert.True(t, len(failRules) == 0, "Some rules are failed:\n%v", failRules)
		assert.Equal(t, len(rules), expectedRuleCount[set.ID()])
	}
}
