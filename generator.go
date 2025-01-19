package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/HomayoonAlimohammadi/structgen/parser"
)

func main() {
	flg := flags{}
	if err := flg.parse(); err != nil {
		log.Fatalf("Failed to parse flags: %v\n", err)
	}

	inputFilesPaths := strings.Split(flg.inputFilesStr, ",")
	if len(inputFilesPaths) == 0 {
		log.Fatalf("No input files provided\n")
	}

	if _, err := os.Stat(flg.outDir); os.IsNotExist(err) {
		if err := os.Mkdir(flg.outDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v\n", err)
		}
	}

	parser, err := parser.New(
		parser.FileType(flg.inputType),
		parser.ParserOptions{
			AdvancedTypesEnabled: flg.advancedTypesEnabled,
			GenerateCmd:          fmt.Sprintf("./%s %s", toolName, strings.Join(os.Args[1:], " ")),
			PkgName:              flg.pkgName,
			OutputDir:            flg.outDir,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create parser: %v\n", err)
	}

	for _, inputFilePath := range inputFilesPaths {
		rcp, err := parser.Parse(inputFilePath)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %v\n", inputFilePath, err)
		}

		if err := rcp.GenerateGoFile(templateFilePath); err != nil {
			log.Fatalf("Failed to generate Go file: %v\n", err)
		}
	}
}
