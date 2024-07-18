package library

import (
	"github.com/dozer111/projectlinter-core/rules"
	file "github.com/dozer111/projectlinter-core/rules/file/rule"
	composerConfig "github.com/dozer111/projectlinter-core/rules/php/composer/config"
	"github.com/dozer111/projectlinter-core/rules/php/composer/config/composer_json"
	"github.com/dozer111/projectlinter-core/rules/php/composer/parser"
	"github.com/dozer111/projectlinter-core/rules/php/composer/rule"
	"github.com/dozer111/projectlinter-core/rules/php/composer/rule/config/platform"
	"github.com/dozer111/projectlinter-core/util/path_provider"
	"os"
)

type composerSet struct {
	pathProvider *path_provider.PathProvider
	config       *composerConfig.Config
	rootEntries  []os.DirEntry

	// additional data for makefileSet.Init
	initialized bool
	initErrors  []error
}

var _ rules.Set = (*composerSet)(nil)

func NewComposerSet(
	pathProvider *path_provider.PathProvider,
	rootEntries []os.DirEntry,
) *composerSet {
	return &composerSet{
		pathProvider: pathProvider,
		rootEntries:  rootEntries,
	}
}

func (s *composerSet) ID() string {
	return "composer"
}

func (s *composerSet) Init() []error {
	s.initialized = true
	composerParser := parser.NewParser(s.pathProvider.PathToCaller())
	composerJson, composerLock, err := composerParser.Parse()

	if err != nil {
		s.initErrors = []error{err}
		return s.initErrors
	}

	s.config = composerConfig.NewComposerConfig(composerJson, composerLock)

	return nil
}

func (s *composerSet) Initialized() bool {
	return s.initialized
}

func (s *composerSet) InitErrors() []error {
	return s.initErrors
}

func (s *composerSet) Config() *composerConfig.Config {
	return s.config
}

func (s *composerSet) Run() *rules.RuleTree {
	r := s.configureRules()

	return rules.NewRuleTree(
		leaf(r.preferStableIsAbsent),
		leaf(r.minimumStabilityIsAbsent),
		leaf(
			r.typeExists,
			leaf(r.typeIsCorrect),
		),
		leaf(
			r.licenseExists,
			leaf(r.licenceIsCorrect),
		),
		leaf(
			r.configExists,
			leaf(
				r.configSortPackagesExists,
				leaf(r.configSortPackagesIsCorrect),
			),
			leaf(r.configPlatformIsAbsent),
		),
		leaf(r.dependencyPHPExists),
		leaf(r.dependencyPHPCsFixerExists),
		leafWithConditions(
			r.phpCsFixerConfigfileIsLatest,
			r.dependencyPHPCsFixerExists.IsPassed,
		),
		leaf(r.dependencyRectorExists),
		leaf(r.descriptionExists),
		leaf(r.dependencyExtJsonIsAbsent),
	)
}

type composerRules struct {
	preferStableIsAbsent         rules.Rule
	minimumStabilityIsAbsent     rules.Rule
	typeExists                   rules.Rule
	typeIsCorrect                rules.Rule
	licenseExists                rules.Rule
	licenceIsCorrect             rules.Rule
	descriptionExists            rules.Rule
	configExists                 rules.Rule
	configSortPackagesExists     rules.Rule
	configSortPackagesIsCorrect  rules.Rule
	configPlatformIsAbsent       rules.Rule
	dependencyPHPExists          rules.Rule
	dependencyPHPCsFixerExists   rules.Rule
	phpCsFixerConfigfileIsLatest rules.Rule
	dependencyRectorExists       rules.Rule
	dependencyExtJsonIsAbsent    rules.Rule
}

/*
configureRules

why is everything done this way, and not just configured in Run?

Because there are a lot of rules in this ruleSet. The most important thing I tried to achieve in Run was simplicity and readability.
To make it convenient for you to understand what is the rules hierarchy - I move rules initialization to another place

This is the philosophy of my approach: many rules - configure separately, keep Run as small and clear as possible
There are few rules - readability will not be greatly affected, you can instantiate directly in Run
*/
func (s *composerSet) configureRules() composerRules {
	const (
		// https://your_github.com/projects/your_php_library_template/browse/composer.json#3
		expectedType = "library"
		// https://your_github.com/projects/your_php_library_template/browse/composer.json#5
		expectedLicence            = "proprietary"
		expectedConfigSortPackages = true
	)

	var (
		configSortPackages    *bool
		configSortPackagesVal bool
		configPlatform        *map[string]string
	)

	if s.config.Config != nil {
		configSortPackages = s.config.Config.SortPackages
		if configSortPackages != nil {
			configSortPackagesVal = *configSortPackages
		}

		configPlatform = s.config.Config.Platform
	}

	boolTrue := true

	return composerRules{
		preferStableIsAbsent:     rule.NewSectionIsAbsentRule("prefer-stable", s.config.PreferStable),
		minimumStabilityIsAbsent: rule.NewSectionIsAbsentRule("minimum-stability", s.config.MinimumStability),
		/*
			Unification. When run `composer bump` without this section will throw a warning:
			Warning: Bumping dependency constraints is not recommended for libraries as it will narrow down your dependencies and may cause problems for your users.
			If your package is not a library, you can explicitly specify the "type" by using "composer config type project".
			Alternatively you can use --dev-only to only bump dependencies within "require-dev".

			https://your_github.com/projects/your_php_library_template/browse/composer.json#3
		*/
		typeExists:    rule.NewSectionExistsRule("type", s.config.Type, expectedType),
		typeIsCorrect: rule.NewSectionHasCorrectValueRule("type", expectedType, *s.config.Type),
		// Composer validate will return warning: `- No license specified, it is recommended to do so.`
		// https://your_github.com/projects/your_php_library_template/browse/composer.json#5
		licenseExists:    rule.NewSectionExistsRule("license", s.config.Licence, expectedLicence),
		licenceIsCorrect: rule.NewSectionHasCorrectValueRule("license", expectedLicence, *s.config.Licence),
		// "composer validate" will return error: `- description : The property description is required`,
		descriptionExists: rule.NewSectionExistsRule("description", s.config.Description, ""),
		configExists: rule.NewSectionExistsRule("config", s.config.Config, composer_json.RawComposerJsonConfigSection{
			SortPackages: &boolTrue,
		}),
		// I can't imagine the case then it must be false,
		// but, why not?
		configSortPackagesExists:    rule.NewSectionExistsRule("config:sort-packages", configSortPackages, expectedConfigSortPackages),
		configSortPackagesIsCorrect: rule.NewSectionHasCorrectValueRule("config:sort-packages", expectedConfigSortPackages, configSortPackagesVal),
		configPlatformIsAbsent:      platform.NewConfigPlatformIsAbsentRule(configPlatform),
		// all the unitedCore apps are with PHP >=8.0
		dependencyPHPExists: rule.NewSpecialDependencyExistsRule(
			s.config.Dependencies,
			"php",
			"^8.0",
			false,
		),
		dependencyPHPCsFixerExists: rule.NewComposerDependencyExistsRule(
			s.config.Dependencies,
			"friendsofphp/php-cs-fixer",
			true,
			[]string{
				"https://your_wiki.com/how-to-wokr-with-phpcsFixer",
				"https://your_git.com/php-library-template/php-cs-fixer.dist",
			},
		),
		phpCsFixerConfigfileIsLatest: file.NewSubstituteFileRule(
			s.rootEntries,
			".php_cs.dist",
			".php-cs-fixer.dist.php",
			"https://your_github.com/projects/your_php_library_template/browse/.php-cs-fixer.dist.php",
			[]string{},
		),
		dependencyRectorExists: rule.NewComposerDependencyExistsRule(
			s.config.Dependencies,
			"rector/rector",
			true,
			[]string{
				"https://your_wiki.com/hot-to-work-with-rector",
				"https://your_git.com/php-library-template/rector.php",
			},
		),
		dependencyExtJsonIsAbsent: rule.NewComposerDependencyIsAbsentRule(
			s.config.Dependencies,
			"ext-json",
			nil,
		),
	}
}
