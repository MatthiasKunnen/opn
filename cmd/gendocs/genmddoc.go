package main

import (
	"fmt"
	"github.com/MatthiasKunnen/opn/cmd/opn"
	"github.com/spf13/cobra/doc"
	"log"
	"os"
	"path"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	outputDir := path.Join(workingDir, "cli_docs")

	err = os.RemoveAll(outputDir)
	if err != nil {
		log.Fatalf("Failed to remove output directory \"%s\": %v", outputDir, err)
	}

	err = os.Mkdir(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory \"%s\": %v", outputDir, err)
	}

	err = doc.GenMarkdownTree(opn.GetCommand(), outputDir)
	if err != nil {
		log.Fatalf("Failed to generate markdown files: %v", err)
	}

	err = os.Symlink("opn.md", path.Join(outputDir, "README.md"))
	if err != nil {
		log.Fatalf("Failed to symlink opn.md to README.md: %v", err)
	}

	fmt.Printf("Generated documentation files in \"%s\".\n", outputDir)
}
