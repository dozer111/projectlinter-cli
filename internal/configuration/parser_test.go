package configuration_test

import (
	"errors"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/configuration"
	utilTest "github.com/dozer111/projectlinter-core/util/test"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigParserSuccessCase(t *testing.T) {
	expectedConfig := &configuration.Configuration{
		"platform",
		"PHP",
		"application",
		[]string{
			"misc",
			"misc.werf_123",
			"phpunit",
		},
	}

	parser := configuration.NewParser(
		utilTest.PathInProjectLinter("testdata/success/" + configuration.ConfigFileName),
	)
	actualConfig, err := parser.Parse()

	assert.Nil(t, err)
	assert.True(t, reflect.DeepEqual(expectedConfig, actualConfig))
}

func TestConfigParserReturnErrorIfConfigIsAbsent(t *testing.T) {
	parser := configuration.NewParser(
		utilTest.PathInProjectLinter("testdata/config_absent/" + configuration.ConfigFileName),
	)
	actualConfig, err := parser.Parse()

	assert.Error(t, err)
	assert.True(t, errors.Is(err, configuration.ConfigIsAbsent))
	assert.Nil(t, actualConfig)
}

func TestConfigParserReturnErrorIfConfigDoesNotApplySchema(t *testing.T) {
	cases := []string{
		"language_absent",
		"language_is_not_one_of_enum",
		"mode_absent",
		"mode_is_not_one_of_enum",
		"unit_absent",
	}

	for _, tc := range cases {
		t.Run(strings.ReplaceAll(tc, "_", " "), func(t *testing.T) {
			parser := configuration.NewParser(
				utilTest.PathInProjectLinter("testdata/config_does_not_apply_schema/" + tc + "/" + configuration.ConfigFileName),
			)
			actualConfig, err := parser.Parse()

			assert.Error(t, err)
			assert.True(t, errors.Is(err, configuration.ConfigDoesNotApplyJsonSchema))
			assert.Nil(t, actualConfig)
		})
	}

}
