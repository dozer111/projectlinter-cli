package application

import (
	"embed"
	"errors"
	utilFile "github.com/dozer111/projectlinter-cli/app-layout/pkg/file"
	"github.com/dozer111/projectlinter-core/rules"
	file "github.com/dozer111/projectlinter-core/rules/file/rule"
	composerConfig "github.com/dozer111/projectlinter-core/rules/php/composer/config"
	"github.com/dozer111/projectlinter-core/rules/php/composer/config/composer_json"
	"github.com/dozer111/projectlinter-core/rules/php/composer/rule/scripts"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"io/fs"
)

const makefileFileName = "Makefile"

type makefileSet struct {
	pathProvider *path_provider.PathProvider
	sources      embed.FS
	embeddedFile fs.File

	composerSet    *composerSet
	composerConfig *composerConfig.Config
}

var _ rules.Set = (*makefileSet)(nil)

func NewMakefileSet(
	pathProvider *path_provider.PathProvider,
	sources embed.FS,
	composerSet *composerSet,
) *makefileSet {
	return &makefileSet{
		pathProvider: pathProvider,
		sources:      sources,
		composerSet:  composerSet,
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

	if !s.composerSet.Initialized() {
		return []error{errors.New("composer set is not initialized")}
	}

	if len(s.composerSet.InitErrors()) > 0 {
		return []error{errors.New("composer set initialize with errors")}
	}

	s.composerConfig = s.composerSet.Config()

	return nil
}

func (s *makefileSet) Run() *rules.RuleTree {
	r := s.configureRules()

	return rules.NewRuleTree(
		// https://your_github.com/projects/your_php_template/browse/Makefile
		leaf(
			r.fileExists,
			leaf(
				r.fileIsLatest,
				// previously we used values-local.yaml to up the aws locally (via minikube)
				// then the targets for this were removed in the Makefile,
				// so now it is not possible to up the aws locally via a new Makefile
				leaf(r.valuesLocalIsAbsent),
				// we changed the file name values--aws-prod.yaml to values-prod.yaml to be the same as in GO application template
				// https://your_github.com/projects/your_php_template/commits/9932b38649f3ead2610524ce1ea7a1783e510ac6#.helm/values-prod.yaml
				leaf(r.valuesProdHasCorrectName),
				// same for secret-values-aws-prod.yaml
				leaf(r.secretValuesProdHasCorrectName),
				// https://your_github.com/projects/your_php_template/pull-requests/86/overview
				// here we added Makefile targets to run the linter only on changed files
				// in order for it to work, composer must additionally have a number of scripts
				leaf(
					r.composerHasScriptsSection,
					leaf(r.composerHasRectorScript),
					leaf(r.composerHasPHPCsFixerScript),
				),
			),
		),
	)
}

type makefileRules struct {
	fileExists                     rules.Rule
	fileIsLatest                   rules.Rule
	valuesLocalIsAbsent            rules.Rule
	valuesProdHasCorrectName       rules.Rule
	secretValuesProdHasCorrectName rules.Rule

	composerHasScriptsSection   rules.Rule
	composerHasRectorScript     rules.Rule
	composerHasPHPCsFixerScript rules.Rule
}

/*
configureRules

why is everything done this way, and not just configured in Run?

Because there are a lot of rules in this ruleSet. The most important thing I tried to achieve in Run was simplicity and readability.
To make it convenient for you to understand what is the rules hierarchy - I move rules initialization to another place

This is the philosophy of my approach: many rules - configure separately, keep Run as small and clear as possible
There are few rules - readability will not be greatly affected, you can instantiate directly in Run
*/
func (s *makefileSet) configureRules() makefileRules {
	const linkOnSource = "https://your_github.com/projects/your_php_template/browse/Makefile"

	composerScriptsSectionValues := s.composerConfig.Scripts
	if composerScriptsSectionValues == nil {
		composerScriptsSectionValues = composer_json.NewScripts(nil, nil, nil)
	}

	return makefileRules{
		fileExists: file.NewFileExistsRule(
			s.pathProvider.PathToCaller(),
			makefileFileName,
			linkOnSource,
		),
		fileIsLatest: file.NewFilesAreSameRule(
			s.pathProvider.PathInCaller(makefileFileName),
			s.embeddedFile,
			linkOnSource,
		),
		valuesLocalIsAbsent: file.NewFileIsAbsentRule(
			s.pathProvider.PathToCaller(),
			".helm/values-local.yaml",
		),
		valuesProdHasCorrectName: file.NewRenameFileRule(
			s.pathProvider.PathToCaller(),
			".helm/values-aws-prod.yaml",
			".helm/values-prod.yaml",
		),
		secretValuesProdHasCorrectName: file.NewRenameFileRule(
			s.pathProvider.PathToCaller(),
			".helm/secret-values-aws-prod.yaml",
			".helm/secret-values-prod.yaml",
		),
		composerHasScriptsSection: scripts.NewScriptsExistsRule(
			s.composerConfig.Scripts,
			composer_json.RawScripts{
				map[string][]string{
					"post-install-cmd": {"@auto-scripts"},
					"post-update-cmd":  {"@auto-scripts"},
				},
				map[string]map[string]string{
					"auto-scripts": {
						"cache:clear": "symfony-cmd",
					},
				},
				map[string]string{
					"php-cs-fixer": "php-cs-fixer fix",
					"rector":       "rector process",
				},
			},
		),
		composerHasRectorScript: scripts.NewScriptsSubsectionExistsRule(
			"rector",
			composerScriptsSectionValues,
			*composer_json.NewScripts(
				nil,
				nil,
				map[string]string{
					"rector": "rector process",
				},
			),
		),
		composerHasPHPCsFixerScript: scripts.NewScriptsSubsectionExistsRule(
			"php-cs-fixer",
			composerScriptsSectionValues,
			*composer_json.NewScripts(
				nil,
				nil,
				map[string]string{
					"php-cs-fixer": "php-cs-fixer fix",
				},
			),
		),
	}
}
