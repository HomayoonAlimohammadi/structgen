package main

import (
	"flag"
	"fmt"
	"slices"

	"github.com/HomayoonAlimohammadi/structgen/parser"
	"golang.org/x/exp/maps"
)

type flags struct {
	inputFilesStr        string
	inputType            string
	pkgName              string
	outDir               string
	advancedTypesEnabled bool
	verbosity            int
}

func (f *flags) parse() error {
	flag.StringVar(&f.inputFilesStr, "files", "", "Comma separated list of input files to generate Go structs from")
	flag.StringVar(&f.pkgName, "pkg", "main", "Name of the package to generate")
	flag.StringVar(&f.outDir, "out-dir", ".", "Directory where the generated files will be saved")
	flag.StringVar(&f.inputType, "type", "", fmt.Sprintf("Type of the input files. Valid types: %s", maps.Keys(parser.ValidTypesToFactory)))
	flag.BoolVar(&f.advancedTypesEnabled, "advanced-types", false, "Enable advanced types (e.g. string instead of any where possible)")
	flag.IntVar(&f.verbosity, "v", 0, "Verbosity level")
	flag.Parse()

	if f.inputType == "" {
		return fmt.Errorf("input type is required, use -type flag")
	}
	if !slices.Contains(maps.Keys(parser.ValidTypesToFactory), parser.FileType(f.inputType)) {
		return fmt.Errorf("invalid input type: %s", f.inputType)
	}

	return nil
}
