package parser

import (
	"fmt"

	"github.com/HomayoonAlimohammadi/structgen/parser/iface"
)

type FileType string

const (
	YAML FileType = "yaml"
	JSON FileType = "json"
)

// ParserOptions represents options for the parser.
type ParserOptions struct {
	// AdvancedTypesEnabled is true if advanced types are enabled.
	AdvancedTypesEnabled bool
	// GenerateCmd is the command that generated the file.
	GenerateCmd string
	// PkgName is the name of the package.
	PkgName string
	// OutputDir is the directory where the generated files will be saved.
	OutputDir string
}

var (
	ValidTypesToFactory = map[FileType]func(opts ParserOptions) (iface.Parser, error){
		YAML: yamlFactory,
		// JSON: jsonFactory,
	}
)

func New(fileType FileType, opts ParserOptions) (iface.Parser, error) {
	if f, ok := ValidTypesToFactory[fileType]; ok {
		return f(opts)
	}
	return nil, fmt.Errorf("invalid file type: %s", fileType)
}
