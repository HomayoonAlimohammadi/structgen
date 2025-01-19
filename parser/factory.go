package parser

import (
	"github.com/HomayoonAlimohammadi/structgen/parser/iface"
	"github.com/HomayoonAlimohammadi/structgen/parser/yaml"
)

func yamlFactory(opts ParserOptions) (iface.Parser, error) {
	yamlOpts := []yaml.Option{}
	if opts.AdvancedTypesEnabled {
		yamlOpts = append(yamlOpts, yaml.WithAdvancedTypesEnabled())
	}
	if opts.GenerateCmd != "" {
		yamlOpts = append(yamlOpts, yaml.WithGenerateCmd(opts.GenerateCmd))
	}
	if opts.PkgName != "" {
		yamlOpts = append(yamlOpts, yaml.WithPkgName(opts.PkgName))
	}
	if opts.OutputDir != "" {
		yamlOpts = append(yamlOpts, yaml.WithOutputDir(opts.OutputDir))
	}
	return yaml.New(yamlOpts...)
}
