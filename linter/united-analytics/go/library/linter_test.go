package library

import (
	"fmt"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/configuration"
	pkgTest "github.com/dozer111/projectlinter-cli/app-layout/pkg/test"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	utilTest "github.com/dozer111/projectlinter-core/util/test"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinterRunAllTheRules(t *testing.T) {
	projectDir := utilTest.PathInProjectLinter("testdata/perfect_project")
	pathProvider := path_provider.NewPathProvider(projectDir)

	config, err := configuration.NewParser(fmt.Sprintf("%s/%s", projectDir, configuration.ConfigFileName)).Parse()
	require.Nil(t, err)
	assert.Equal(t, config.Unit, pkgTest.UnitName())

	l := NewLinter(pathProvider)
	sets, err := l.Sets()
	require.Nil(t, err)

	expectedRuleCount := map[string]int{
		"dependency":    2,
		".editorconfig": 2,
	}

	for _, set := range sets {
		require.True(t, len(set.Init()) == 0)

		rules := set.Run().Resolve(nil)

		failRules := utilTest.AllSetRulesArePassed(rules)
		assert.True(t, len(failRules) == 0, "Some rules are failed:\n%v", failRules)
		assert.Equal(t, len(rules), expectedRuleCount[set.ID()])
	}
}
