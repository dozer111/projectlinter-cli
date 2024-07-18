package application

import (
	"embed"
	utilFile "github.com/dozer111/projectlinter-cli/app-layout/pkg/file"
	"github.com/dozer111/projectlinter-core/rules"
	file "github.com/dozer111/projectlinter-core/rules/file/rule"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"io/fs"
)

const makefileFileName = "Makefile"

type makefileSet struct {
	pathProvider *path_provider.PathProvider

	sources      embed.FS
	embeddedFile fs.File
}

var _ rules.Set = (*makefileSet)(nil)

func NewMakefileSet(
	pathProvider *path_provider.PathProvider,
	sources embed.FS,
) *makefileSet {
	return &makefileSet{
		pathProvider: pathProvider,
		sources:      sources,
	}
}

func (s *makefileSet) ID() string {
	return "makefile"
}

func (s *makefileSet) Init() []error {
	f, err := utilFile.FetchSourceFile(s.sources, makefileFileName)
	if err != nil {
		return []error{err}
	}

	s.embeddedFile = f
	return nil
}

func (s *makefileSet) Run() *rules.RuleTree {
	const linkOnSource = "https://your_github.com/projects/your_template/browse/Makefile"
	fileExists := file.NewFileExistsRule(
		s.pathProvider.PathToCaller(),
		makefileFileName,
		linkOnSource,
	)
	fileIsLatest := file.NewFilesAreSameRule(
		s.pathProvider.PathInCaller(makefileFileName),
		s.embeddedFile,
		linkOnSource,
	)

	valuesProdHasCorrectName := file.NewRenameFileRule(
		s.pathProvider.PathToCaller(),
		".helm/values-aws.yaml",
		".helm/values-prod.yaml",
	)

	secretValuesProdHasCorrectName := file.NewRenameFileRule(
		s.pathProvider.PathToCaller(),
		".helm/secret-values-aws.yaml",
		".helm/secret-values-prod.yaml",
	)

	return rules.NewRuleTree(
		leaf(
			fileExists,
			leaf(
				fileIsLatest,
				leaf(valuesProdHasCorrectName),
				leaf(secretValuesProdHasCorrectName),
			),
		),
	)
}
