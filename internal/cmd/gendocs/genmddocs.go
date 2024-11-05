package gendocs

import (
	"fmt"
	"github.com/MatthiasKunnen/opn/internal/cmd/opn"
	"github.com/spf13/cobra/doc"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// GenMdDocs will remove all markdown files in the dir, create it if necessary, and generate the
// markdown documentation in it.
// outputDir must be an absolute path.
func GenMdDocs(outputDir string) error {
	if !filepath.IsAbs(outputDir) {
		return fmt.Errorf("GenMdDocs: output directory must be an absolute path")
	}

	if outputDir == "/" {
		return fmt.Errorf("GenMdDocs: output directory can't be /. Use a subdir")
	}

	dirEntries, err := os.ReadDir(outputDir)
	if err != nil {
		return fmt.Errorf("GenMdDocs: error reading output directory: %w", err)
	}

	for _, entry := range dirEntries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".md") {
			continue
		}

		name = path.Join(outputDir, name)
		if err := os.Remove(name); err != nil {
			return fmt.Errorf("GenMdDocs: failed to remove existing markdown file: %w", err)
		}
	}

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("GenMdDocs: failed to create output directory: %v", err)
	}

	err = doc.GenMarkdownTree(opn.GetCommand(), outputDir)
	if err != nil {
		return fmt.Errorf("GenMdDocs: failed to generate markdown files: %v", err)
	}

	err = os.Symlink("opn.md", path.Join(outputDir, "README.md"))
	if err != nil {
		return fmt.Errorf("GenMdDocs: failed to symlink opn.md to README.md: %v", err)
	}

	return nil
}
