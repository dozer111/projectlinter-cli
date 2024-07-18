package common_sets

import (
	"embed"
	utilFile "github.com/dozer111/projectlinter-cli/app-layout/pkg/file"
	"github.com/dozer111/projectlinter-core/rules"
	file "github.com/dozer111/projectlinter-core/rules/file/rule"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"io/fs"
)

const editorconfigFileName = ".editorconfig"

type editorconfigSet struct {
	sources      embed.FS
	pathProvider *path_provider.PathProvider

	embeddedFile        fs.File
	linkOnCorrectSource string
}

var _ rules.Set = (*editorconfigSet)(nil)

func NewEditorConfig(
	pathProvider *path_provider.PathProvider,
	sources embed.FS,
	linkOnCorrectSource string,
) *editorconfigSet {
	return &editorconfigSet{
		pathProvider:        pathProvider,
		sources:             sources,
		linkOnCorrectSource: linkOnCorrectSource,
	}
}

func (s *editorconfigSet) ID() string {
	return editorconfigFileName
}

func (s *editorconfigSet) Init() []error {
	f, err := utilFile.FetchSourceFile(s.sources, editorconfigFileName)
	if err != nil {
		return []error{err}
	}

	s.embeddedFile = f

	return nil
}

func (s *editorconfigSet) Run() *rules.RuleTree {
	fileExists := file.NewFileExistsRule(
		s.pathProvider.PathToCaller(),
		editorconfigFileName,
		s.linkOnCorrectSource,
	)

	fileIsLatest := file.NewFilesAreSameRule(
		s.pathProvider.PathInCaller(editorconfigFileName),
		s.embeddedFile,
		s.linkOnCorrectSource,
	)

	return rules.NewRuleTree(
		leaf(
			fileExists,
			leaf(fileIsLatest),
		),
	)
}
