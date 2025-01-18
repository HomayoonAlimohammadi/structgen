package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/HomayoonAlimohammadi/structgen/parser"
	"github.com/HomayoonAlimohammadi/structgen/parser/iface"
	"golang.org/x/exp/maps"
)

const (
	rootStructDocStringFmt = "// %s represents the values of the %s chart"
	templateFilePath       = "struct.go.tmpl"
	toolName               = "CHART_VALUES_STRUCT_GENERATOR"
)

func parseFlags(inputFilesStr, inputType, pkgName, outDir *string, advancedTypesEnabled *bool, verbosity *int) error {
	flag.StringVar(inputFilesStr, "files", "", "Comma separated list of input files to generate Go structs from")
	flag.StringVar(pkgName, "pkg", "main", "Name of the package to generate")
	flag.StringVar(outDir, "out-dir", ".", "Directory where the generated files will be saved")
	flag.StringVar(inputType, "type", "", fmt.Sprintf("Type of the input files. Valid types: %s", maps.Keys(parser.ValidTypesToFactory)))
	flag.BoolVar(advancedTypesEnabled, "advanced-types", false, "Enable advanced types (e.g. string instead of any where possible)")
	flag.IntVar(verbosity, "v", 0, "Verbosity level")
	flag.Parse()

	if *inputType == "" {
		return fmt.Errorf("input type is required, use -type flag")
	}
	if !slices.Contains(maps.Keys(parser.ValidTypesToFactory), parser.FileType(*inputType)) {
		return fmt.Errorf("invalid input type: %s", *inputType)
	}

	return nil
}

func main() {
	var (
		inputFilesStr        string
		inputType            string
		pkgName              string
		outDir               string
		advancedTypesEnabled bool
		verbosity            int
	)

	parseFlags(&inputFilesStr, &inputType, &pkgName, &outDir, &advancedTypesEnabled, &verbosity)

	inputFilesPaths := strings.Split(inputFilesStr, ",")
	if len(inputFilesPaths) == 0 {
		log.Fatalf("No input files provided\n")
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.Mkdir(outDir, 0755)
		if err != nil {
			log.Fatalf("Failed to create output directory: %v\n", err)
		}
	}

	parser, err := parser.New(parser.FileType(inputType))
	if err != nil {
		log.Fatalf("Failed to create parser: %v\n", err)
	}

	generateCmd := fmt.Sprintf("./%s %s", toolName, strings.Join(os.Args[1:], " "))

	parseOpts := iface.ParseOptions{
		AdvancedTypesEnabled: advancedTypesEnabled,
		GenerateCmd:          generateCmd,
	}
	for _, inputFilePath := range inputFilesPaths {
		rcp, err := parser.Parse(inputFilePath, outDir, pkgName, parseOpts)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %v\n", inputFilePath, err)
		}

		if err := rcp.GenerateGoFile(templateFilePath); err != nil {
			log.Fatalf("Failed to generate Go file: %v\n", err)
		}
	}
}
