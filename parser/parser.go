package parser

import (
	"fmt"

	"github.com/HomayoonAlimohammadi/structgen/parser/iface"
	"github.com/HomayoonAlimohammadi/structgen/parser/yaml"
)

type FileType string

const (
	YAML FileType = "yaml"
	JSON FileType = "json"
)

var (
	ValidTypesToFactory = map[FileType]func() (iface.Parser, error){
		YAML: yaml.New,
		// JSON: json.New,
	}
)

func New(fileType FileType) (iface.Parser, error) {
	if f, ok := ValidTypesToFactory[fileType]; ok {
		return f()
	}
	return nil, fmt.Errorf("invalid file type: %s", fileType)
}
