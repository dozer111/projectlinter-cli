package configuration

import (
	_ "embed"
	"errors"
	"fmt"
	"os"

	"github.com/1set/gut/yos"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed json-schema.json
var jsonSchema string

type Parser struct {
	pathToConfig string
}

func NewParser(pathToConfig string) *Parser {
	return &Parser{pathToConfig}
}

var (
	ConfigIsAbsent               = fmt.Errorf("configuration file is absent")
	ConfigDoesNotApplyJsonSchema = fmt.Errorf("configuration file does not satisfy its json-schema")
)

func (p *Parser) Parse() (*Configuration, error) {
	if !yos.ExistFile(p.pathToConfig) {
		return nil, ConfigIsAbsent
	}

	data, err := os.ReadFile(p.pathToConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", p.pathToConfig, err)
	}

	if err := p.assertConfigSatisfiesToJsonSchema(data, p.pathToConfig); err != nil {
		return nil, err
	}

	var cfg Configuration
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("Configuration.Unmarshal failed: %w", err)
	}

	return &cfg, nil
}

func (p *Parser) assertConfigSatisfiesToJsonSchema(data []byte, filePath string) error {
	schemaLoader := gojsonschema.NewStringLoader(jsonSchema)
	var rawConfig interface{}
	err := yaml.Unmarshal(data, &rawConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot parse bytes from %s to interface{}: %v", filePath, err))
	}
	documentLoader := gojsonschema.NewRawLoader(rawConfig)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot check does the file %s satisfy json-schema: %v", filePath, err))
	}
	if result.Valid() == false {
		return fmt.Errorf("%w: %s: %v", ConfigDoesNotApplyJsonSchema, filePath, result.Errors())
	}

	return nil
}
