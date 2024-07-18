package common_sets_test

import (
	"embed"
	"github.com/dozer111/projectlinter-cli/app-layout/internal/common_sets"
	utilFile "github.com/dozer111/projectlinter-cli/app-layout/pkg/file"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	utilTest "github.com/dozer111/projectlinter-core/util/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/source/*
var editorconfig embed.FS

func TestEditorconfigSet(t *testing.T) {
	utilFile.FetchSourcePath = "testdata/source/"
	pathProvider := path_provider.NewPathProvider(utilTest.PathInProjectLinter("testdata/editorconfig"))

	s := common_sets.NewEditorConfig(pathProvider, editorconfig, "")
	errs := s.Init()
	assert.Equal(t, 0, len(errs))

	rules := s.Run().Resolve([]string{})

	assert.Equal(t, 2, len(rules))
	failRules := utilTest.AllSetRulesArePassed(rules)
	assert.True(t, len(failRules) == 0, "Some rules are failed:\n%v", failRules)
}
