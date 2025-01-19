package iface

import (
	"github.com/HomayoonAlimohammadi/structgen/parser/recipe"
)

// Parser is an interface for parsing a file.
type Parser interface {
	Parse(filePath string) (*recipe.StructsRecipe, error)
}
