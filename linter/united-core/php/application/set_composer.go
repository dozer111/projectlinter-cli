package application

import (
	"github.com/Masterminds/semver/v3"
	pkg "github.com/dozer111/projectlinter-cli/app-layout/pkg/generic"
	"github.com/dozer111/projectlinter-core/rules"
	checkFile "github.com/dozer111/projectlinter-core/rules/file/rule"
	composerConfig "github.com/dozer111/projectlinter-core/rules/php/composer/config"
	"github.com/dozer111/projectlinter-core/rules/php/composer/config/composer_json"
	composerParser "github.com/dozer111/projectlinter-core/rules/php/composer/parser"
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
	p := composerParser.NewParser(s.pathProvider.PathToCaller())
	composerJson, composerLock, err := p.Parse()

	if err != nil {
		s.initErrors = []error{err}
		return s.initErrors
	}

	s.config = composerConfig.NewComposerConfig(composerJson, composerLock)
	return nil
}

func (s *composerSet) InitErrors() []error {
	return s.initErrors
}

func (s *composerSet) Initialized() bool {
	return s.initialized
}

func (s *composerSet) Config() *composerConfig.Config {
	return s.config
}

func (s *composerSet) Run() *rules.RuleTree {
	r := s.configureRules()

	return rules.NewRuleTree(
		leaf(r.preferStableIsAbsent),
		leaf(r.minimumStabilityIsAbsent),
		/*
			Unification. When run `composer bump` without this section will throw a warning:
			Warning: Bumping dependency constraints is not recommended for libraries as it will narrow down your dependencies and may cause problems for your users.
			If your package is not a library, you can explicitly specify the "type" by using "composer config type project".
			Alternatively you can use --dev-only to only bump dependencies within "require-dev".

			https://your_github.com/projects/your_php_template/browse/composer.json#5
		*/
		leaf(
			r.typeExists,
			leaf(r.typeIsCorrect),
		),

		// Composer validate will return warning: `- No license specified, it is recommended to do so.`
		// https://your_github.com/projects/your_php_template/browse/composer.json#4
		leaf(
			r.licenseExists,
			leaf(r.licenceIsCorrect),
		),
		// PHP dependency is REQUIRED!
		leaf(r.dependencyPHPExists),
		leaf(
			r.configExists,
			leaf(
				r.configSortPackagesExists,
				leaf(r.configSortPackagesIsCorrect),
			),
			// https://your_github.com/projects/your_php_template/browse/composer.json#64-66
			// This setting is CRUTIAL for application
			// We have discuss it with topTeamlead1, topTeamlead2, topColleague3 and agreed that this setting MUST BE!
			rules.NewLeafWithConditions(
				r.configPlatformExists,
				[]func() bool{r.dependencyPHPExists.IsPassed},
				leaf(r.configPlatformHasPHP),
				// currently 2023.12.04 I don't know of any other platform that should still be in the rules.
				leaf(r.configPlatformContainOnlyPHP),
			),
		),
		leaf(r.dependencyPHPCsFixerExists),
		leafWithConditions(
			r.phpCsFixerConfigfileIsLatest,
			r.dependencyPHPCsFixerExists.IsPassed,
		),
		leaf(r.dependencyRectorExists),
		// "composer validate" will return error: `- description : The property description is required`,
		leaf(r.descriptionExists),
		leaf(r.dependenciesConstraintsAreValid),
		leaf(r.dependencyExtJsonIsAbsent),
		// Problem: reading POST request body
		// before symfony 6.3, it was impossible to get easy assocArray from requestBody
		// That's why we used the temporary solution "qandidate/symfony-json-request-transformer" dependency
		//
		// With the advent of symfony6.3, we no longer need this dependency
		leaf(
			r.dependencySymfonyFrameworkBundleExists,
			leafWithConditions(
				r.dependencyQuandidateJsonTransformerIsAbsent,
				func() bool {
					symfonyDependency := s.config.Dependencies.Get(symfony)
					return symfonyDependency.VersionIsCorrect() && symfonyDependency.Version().GreaterThan(semver.MustParse("6.3"))
				},
			),
		),
	)
}

type composerRules struct {
	preferStableIsAbsent                        rules.Rule
	minimumStabilityIsAbsent                    rules.Rule
	typeExists                                  rules.Rule
	typeIsCorrect                               rules.Rule
	licenseExists                               rules.Rule
	licenceIsCorrect                            rules.Rule
	dependencyPHPExists                         rules.Rule
	configExists                                rules.Rule
	configSortPackagesExists                    rules.Rule
	configSortPackagesIsCorrect                 rules.Rule
	configPlatformExists                        rules.Rule
	configPlatformHasPHP                        rules.Rule
	configPlatformContainOnlyPHP                rules.Rule
	dependencyPHPCsFixerExists                  rules.Rule
	phpCsFixerConfigfileIsLatest                rules.Rule
	dependencyRectorExists                      rules.Rule
	descriptionExists                           rules.Rule
	dependenciesConstraintsAreValid             rules.Rule
	dependencyExtJsonIsAbsent                   rules.Rule
	dependencyQuandidateJsonTransformerIsAbsent rules.Rule
	// by this dependency, I am guided by which version of symfony is currently on the project
	dependencySymfonyFrameworkBundleExists rules.Rule
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
		// https://your_github.com/projects/your_php_template/browse/composer.json#5
		expectedType = "project"
		// https://your_github.com/projects/your_php_template/browse/composer.json#4
		expectedLicense      = "proprietary"
		expectedSortPackages = true

		phpDependencyName = "php"
	)

	proposedPlatformPHP := s.config.PHP.TrimConstraint()

	var configPlatform *map[string]string = nil
	var configPlatformValue map[string]string = nil
	var configSortPackages *bool = nil
	var actualSortPackages bool
	if s.config.Config != nil {
		actualSortPackages = pkg.ValueOrDefaultFromPointer(s.config.Config.SortPackages)
		configPlatform = s.config.Config.Platform
		if configPlatform != nil {
			configPlatformValue = *s.config.Config.Platform
		}

		configSortPackages = s.config.Config.SortPackages
	}

	boolTrue := true

	return composerRules{
		preferStableIsAbsent:     rule.NewSectionIsAbsentRule("prefer-stable", s.config.PreferStable),
		minimumStabilityIsAbsent: rule.NewSectionIsAbsentRule("minimum-stability", s.config.MinimumStability),
		typeExists:               rule.NewSectionExistsRule("type", s.config.Type, expectedType),
		typeIsCorrect:            rule.NewSectionHasCorrectValueRule("type", expectedType, pkg.ValueOrDefaultFromPointer(s.config.Type)),
		licenseExists:            rule.NewSectionExistsRule("license", s.config.Licence, expectedLicense),
		licenceIsCorrect:         rule.NewSectionHasCorrectValueRule("license", expectedLicense, pkg.ValueOrDefaultFromPointer(s.config.Licence)),
		// UnitedCore has no service with PHP <8.0
		dependencyPHPExists: rule.NewSpecialDependencyExistsRule(
			s.config.Dependencies,
			phpDependencyName,
			"^8.0",
			false,
		),
		configExists: rule.NewSectionExistsRule(
			"config",
			s.config.Config,
			composer_json.RawComposerJsonConfigSection{SortPackages: &boolTrue},
		),
		configSortPackagesExists: rule.NewSectionExistsRule("config:sort-packages", configSortPackages, expectedSortPackages),
		configSortPackagesIsCorrect: rule.NewSectionHasCorrectValueRule(
			"config:sort-packages",
			expectedSortPackages,
			actualSortPackages,
		),
		configPlatformExists: platform.NewConfigPlatformExistsRule(
			configPlatform,
			map[string]string{phpDependencyName: proposedPlatformPHP},
		),
		configPlatformHasPHP: platform.NewSpecifiedPlatformExistsRule(
			configPlatformValue,
			phpDependencyName,
			proposedPlatformPHP,
		),
		configPlatformContainOnlyPHP: platform.NewOnlySpecifiedPlatformExistsRule(
			configPlatformValue,
			map[string]string{phpDependencyName: proposedPlatformPHP},
		),
		dependencyPHPCsFixerExists: rule.NewComposerDependencyExistsRule(
			s.config.Dependencies,
			"friendsofphp/php-cs-fixer",
			true,
			[]string{
				"https://your_wiki.com/how-to-work-with-phpcsfixer",
				"Copy config from https://your_git.com/projects/your_php_application_template/php-cs-fixer.dist.php",
			},
		),
		phpCsFixerConfigfileIsLatest: checkFile.NewSubstituteFileRule(
			s.rootEntries,
			".php_cs.dist",
			".php-cs-fixer.dist.php",
			"https://your_github.com/projects/your_php_template/browse/.php-cs-fixer.dist.php",
			[]string{},
		),
		dependencyRectorExists: rule.NewComposerDependencyExistsRule(
			s.config.Dependencies,
			"rector/rector",
			true,
			[]string{
				"https://your_wiki.com/how-to-rector",
				"Copy config from https://your_git.com/projects/your_php_application_template/rector.php",
			},
		),
		descriptionExists: rule.NewSectionExistsRule("description", s.config.Description, ""),
		dependenciesConstraintsAreValid: rule.NewDependenciesConstrainsAreValidRule(
			s.config.Dependencies,
			[]string{
				"ext-PDO",
				"ext-apcu",
				"ext-bcmath",
				"ext-curl",
				"ext-ds",
				"ext-fileinfo",
				"ext-iconv",
				"ext-intl",
				"ext-json",
				"ext-libxml",
				"ext-openssl",
				"ext-pcntl",
				"ext-pdo",
				"ext-redis",
				"ext-simplexml",
				"ext-soap",
				"ext-sockets",
				"ext-xmlreader",
				"ext-zlib",
			},
		),
		dependencyExtJsonIsAbsent: rule.NewComposerDependencyIsAbsentRule(
			s.config.Dependencies,
			"ext-json",
			nil,
		),
		dependencyQuandidateJsonTransformerIsAbsent: rule.NewComposerDependencyIsAbsentRule(
			s.config.Dependencies,
			"qandidate/symfony-json-request-transformer",
			[]string{
				"Remove dependency is not enough.",
				"You need to check files with Symfony\\Component\\HttpFoundation\\Request",
				"Change `$request->request->...` on `$request->getPayload()->...`",
				"Docs: https://symfony.com/blog/new-in-symfony-6-3-request-payload",
				"Example(service): https://your_git.com/auth-sv/commits/e5ce2aef39c0403e6b94757cb6ad8d0198b960ea#src/Service/Symfony/Controller/SomeMyController.php",
				"Example(library): https://your_git.com/auth-lb/commits/96ce5f6514fa6f58f43c4b8b755d4048c78b8c09#src/Symfony/Controllers/SomeLBController.php",
			},
		),
		dependencySymfonyFrameworkBundleExists: rule.NewComposerDependencyExistsRule(
			s.config.Dependencies,
			symfony,
			false,
			[]string{
				"Project linter use this dependency to understand what symfony version is on the project",
				"It is very strange that it is absent.",
				"00-base-tpl has it: https://your_github.com/projects/your_php_template/browse/composer.json#25",
			},
		),
	}
}

const symfony = "symfony/framework-bundle"
