package yaml

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/HomayoonAlimohammadi/structgen/parser/iface"
	"github.com/HomayoonAlimohammadi/structgen/parser/recipe"
	"github.com/iancoleman/strcase"
	"gopkg.in/yaml.v3"
)

const (
	toolName               = "CHART_VALUES_STRUCT_GENERATOR"
	rootStructDocStringFmt = "// %s represents the values of the %s chart"
)

type Parser struct {
	advancedTypesEnabled bool
}

func New() (iface.Parser, error) {
	// p := &Parser{}
	// for _, opt := range opts {
	// 	if err := opt.apply(p); err != nil {
	// 		return nil, fmt.Errorf("failed to apply option: %w", err)
	// 	}
	// }

	return &Parser{}, nil
}

func (p *Parser) Parse(filePath, outDir, pkgName string, opts ...iface.ParseOptions) (*recipe.StructsRecipe, error) {
	var opt iface.ParseOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	p.handleParseOptions(opt)

	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file %s does not exist", filePath)
		} else {
			return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
		}
	}
	defer file.Close()

	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	rootStructName := strcase.ToCamel(strings.ReplaceAll(strings.ReplaceAll(baseName, ".", "_"), "-", "_"))
	outputFilePath := path.Join(outDir, fmt.Sprintf("%s.go", baseName))

	rcp := &recipe.StructsRecipe{
		OutputFilePath: outputFilePath,
		RootStructName: rootStructName,
		PkgName:        pkgName,
		GenerateCmd:    opt.GenerateCmd,
		GenerateDate:   time.Now().Format(time.DateOnly),
		ToolName:       toolName,
		Imports: []recipe.Import{
			{
				Path: "fmt",
			},
			{
				Path: "encoding/json",
			},
			{
				Path: "reflect",
			},
			{
				Path: "strings",
			},
		},
	}

	rootNode := yaml.Node{}
	if err := yaml.NewDecoder(file).Decode(&rootNode); err != nil {
		return nil, fmt.Errorf("error decoding yaml value from file %s: %w", filePath, err)
	}

	if len(rootNode.Content) == 0 {
		return nil, fmt.Errorf("empty file %s", filePath)
	}

	docString := fmt.Sprintf(rootStructDocStringFmt, rootStructName, filePath)
	p.parse(rcp, rootStructName, rootNode.Content[0], docString, true)

	return rcp, nil
}

// parse recursively generates Go recipe definitions from a YAML Node
func (p *Parser) parse(rcp *recipe.StructsRecipe, structName string, node *yaml.Node, docString string, isRoot bool) {
	stMeta := &recipe.StructMeta{
		IsRoot:    isRoot,
		Name:      structName,
		DocString: docString,
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]
		fieldName := strcase.ToCamel(keyNode.Value)

		field := &recipe.FieldMeta{
			Name:         fieldName,
			OriginalName: keyNode.Value,
			DocString:    strings.Join(extractComments(keyNode, valueNode), "\n"),
		}

		// TODO: handle such cases:
		// controller:
		// 	<<: *defaults
		if keyNode.Value == "<<" {
			continue
		}

		switch valueNode.Kind {
		case yaml.MappingNode:
			// nested struct
			if len(valueNode.Content) == 0 {
				field.Type = infereTypeString(valueNode, p.advancedTypesEnabled, false)
			} else {
				// struct of known type, the type will be the name of the struct
				nestedStructName := structName + "_" + fieldName
				field.Type = "*" + nestedStructName
				p.parse(rcp, nestedStructName, valueNode, field.DocString, false)
			}
		case yaml.SequenceNode:
			if len(valueNode.Content) == 0 || len(valueNode.Content[0].Content) == 0 {
				field.Type = infereTypeString(valueNode, p.advancedTypesEnabled, false)
			} else {
				// list with its own struct
				nestedListName := structName + "_" + fieldName + "Item"
				field.Type = "*[]" + nestedListName
				p.parse(rcp, nestedListName, valueNode.Content[0], field.DocString, false)
			}
		case yaml.ScalarNode:
			// scalar value
			field.Type = infereTypeString(valueNode, p.advancedTypesEnabled, false)
		}

		stMeta.Fields = append(stMeta.Fields, field)
	}

	rcp.Structs = append(rcp.Structs, stMeta)
}

func (p *Parser) handleParseOptions(opt iface.ParseOptions) {
	p.advancedTypesEnabled = opt.AdvancedTypesEnabled
}

// infereTypeString infers the Go type of a YAML node
func infereTypeString(n *yaml.Node, advanced bool, isNested bool) string {
	switch n.Kind {
	case yaml.ScalarNode:
		if !advanced {
			return "any"
		}

		if n.Value == "~" {
			return "any"
		}

		switch n.Tag {
		case "!!bool":
			if isNested {
				return "bool"
			}
			return "*bool"
		case "!!int":
			if isNested {
				return "int"
			}
			return "*int"
		case "!!float":
			if isNested {
				return "float64"
			}
			return "*float64"
		default:
			if isNested {
				return "string"
			}
			return "*string"
		}
	case yaml.SequenceNode:
		if len(n.Content) == 0 || !advanced {
			return "*[]any"
		}
		return "*[]" + infereTypeString(n.Content[0], true, true) // advanced has to be true
	case yaml.MappingNode:
		// advanced inference for maps should be handled by the upper level
		return "*map[string]any"
	default:
		return "any"
	}
}

// extractComments extracts comments from a YAML node
func extractComments(keyNode, valNode *yaml.Node) []string {
	totalLines := []string{}

	if hc := keyNode.HeadComment; hc != "" {
		lines := strings.Split(hc, "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			l = strings.TrimLeft(l, "#")
			totalLines = append(totalLines, fmt.Sprintf("// %s", l))
		}
	}

	if lc := keyNode.LineComment; lc != "" {
		lines := strings.Split(lc, "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			l = strings.TrimLeft(l, "#")
			totalLines = append(totalLines, fmt.Sprintf("// %s", l))
		}
	}

	if fc := keyNode.FootComment; fc != "" {
		lines := strings.Split(fc, "\n")
		for _, l := range lines {
			l = strings.TrimSpace(l)
			l = strings.TrimLeft(l, "#")
			totalLines = append(totalLines, fmt.Sprintf("// %s", l))
		}
	}

	if valNode.Value != "" {
		if len(totalLines) != 0 {
			totalLines = append(totalLines, "//")
		}
		totalLines = append(totalLines, fmt.Sprintf("// Default value in yaml: %s", valNode.Value))
	} else if valNode.Kind == yaml.SequenceNode {
		if len(valNode.Content) != 0 && len(valNode.Content[0].Content) == 0 {
			if len(totalLines) != 0 {
				totalLines = append(totalLines, "//")
			}
			totalLines = append(totalLines, "// Default value in yaml:")
			for _, c := range valNode.Content {
				totalLines = append(totalLines, fmt.Sprintf("// - %s", c.Value))
			}
		}
	}

	return totalLines
}
