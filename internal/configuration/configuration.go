package configuration

import "slices"

const ConfigFileName = "project-linter.yaml"

type (
	Configuration struct {
		Unit     string   `yaml:"unit"`
		Language Language `yaml:"language"`
		Mode     Mode     `yaml:"mode"`
		Ignore   []string `yaml:"ignore,omitempty"`
	}

	Language string
	Mode     string
)

var AvailableLanguages = []string{
	string(LanguageGO),
	string(LanguagePHP),
}

var AvailableMods = []string{
	string(ModeApplication),
	string(ModeLibrary),
}

const (
	LanguageGO  Language = "GO"
	LanguagePHP Language = "PHP"

	ModeApplication Mode = "application"
	ModeLibrary     Mode = "library"
)

func (l Language) IsGO() bool {
	return l == LanguageGO
}

func (l Language) IsPHP() bool {
	return l == LanguagePHP
}

func (l Language) Valid() bool {
	return slices.Contains[[]string, string](AvailableLanguages, string(l))
}

func (m Mode) IsApplication() bool {
	return m == ModeApplication
}

func (m Mode) Valid() bool {
	return slices.Contains[[]string, string](AvailableMods, string(m))
}
