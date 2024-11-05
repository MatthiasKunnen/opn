package main

import (
	"fmt"
	"github.com/MatthiasKunnen/opn/internal/cmd/gendocs"
	"log"
	"os"
	"path"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	outputDir := path.Join(workingDir, "docs/cli")
	err = gendocs.GenMdDocs(outputDir)
	if err != nil {
		log.Fatalf("Failed to generate documentation: %v", err)
	}

	fmt.Printf("Generated documentation files in \"%s\".\n", outputDir)
}
