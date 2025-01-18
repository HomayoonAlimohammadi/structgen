package iface

import (
	"github.com/HomayoonAlimohammadi/structgen/parser/recipe"
)

// Parser is an interface for parsing a file.
type Parser interface {
	Parse(filePath, outDir, pkgName string, opts ...ParseOptions) (*recipe.StructsRecipe, error)
}

// ParseOptions represents options for parsing a file.
type ParseOptions struct {
	// AdvancedTypesEnabled is true if advanced types are enabled.
	AdvancedTypesEnabled bool
	// GenerateCmd is the command that generated the file.
	GenerateCmd string
}
