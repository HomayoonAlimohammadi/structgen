package recipe

import (
	"fmt"
	"html/template"
	"os"

	"github.com/HomayoonAlimohammadi/structgen/parser/utils"
)

// StructMeta represents a struct in a Go file
type StructMeta struct {
	// IsRoot is true if the struct is the root struct of the file.
	// This struct represents the complete set of values of the input file.
	IsRoot bool
	// Name is the Name of the struct.
	Name string
	// DocString is the docstring of the struct.
	// different lines should be separated by \n.
	DocString string
	// Fields is a list of Fields in the struct.
	Fields []*FieldMeta
}

// FieldMeta represents a field in a Go struct
type FieldMeta struct {
	// Name is the Name of the field.
	Name string
	// OriginalName is the original name of the field in the input file.
	OriginalName string
	// DocString is the docstring of the field.
	// different lines should be separated by \n.
	DocString string
	// Type is the go type of the field.
	Type string
}

// Import represents an import in a Go file
type Import struct {
	// Alias is the alias of the import.
	Alias string
	// Path is the path of the import.
	Path string
}

// StructsRecipe is a recipe to generate a Go file from an input file
type StructsRecipe struct {
	// OutputFilePath is the path of the output file.
	OutputFilePath string
	// RootStructName is the name of the root struct.
	RootStructName string

	// GenerateCmd is the command that generated the file.
	GenerateCmd string
	// GenerateDate is the date the file was generated.
	GenerateDate string
	// ToolName is the name of the tool that generated the file.
	ToolName string

	// PkgName is the name of the package.
	PkgName string
	// Imports is a list of Imports in the file.
	Imports []Import
	// Structs is a list of Structs in the file.
	Structs []*StructMeta
}

// GenerateGoFile generates a Go file from a recipe
func (r *StructsRecipe) GenerateGoFile(templateFilePath string) error {
	tmpl, err := template.New(templateFilePath).ParseFiles(templateFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse template file %s: %w", templateFilePath, err)
	}

	var out *os.File
	out, err = os.Create(r.OutputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create Go file %s: %w", r.OutputFilePath, err)
	}

	if err := out.Chmod(0644); err != nil {
		return fmt.Errorf("failed to change permissions of Go file %s: %w", r.OutputFilePath, err)
	}

	if err := tmpl.Execute(out, r); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := utils.FormatGoFile(r.OutputFilePath); err != nil {
		return fmt.Errorf("failed to format Go file %s: %w", r.OutputFilePath, err)
	}

	return nil
}
