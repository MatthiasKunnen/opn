package opnlib

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os/exec"
	"strings"
)

// GetFileMime returns the mime type of the given path.
// It uses programs installed
func GetFileMime(path string) (string, error) {
	xdgCmd := exec.Command("xdg-mime", "query", "filetype", path)
	output, err := xdgCmd.Output()

	switch {
	case errors.Is(err, exec.ErrNotFound):
		// Fall back to file
	case err != nil:
		log.Printf(
			"Failed to use xdg-mime to establish MIME type, falling back to file: %v\n",
			err,
		)
	default:
		return strings.TrimSpace(string(output)), nil
	}

	fileCmd := exec.Command(
		"/usr/bin/file",
		"--brief",
		"--dereference",
		"--mime-type",
		path,
	)
	output, err = fileCmd.Output()

	switch {
	case errors.Is(err, exec.ErrNotFound), errors.Is(err, fs.ErrNotExist):
		return "", fmt.Errorf("failed to determine MIME type, no programs to determine" +
			" MIME type are installed. Either xdg-mime (xdg-utils) or, file, is required")
	case err != nil:
		return "", err
	default:
		return strings.TrimSpace(string(output)), nil
	}
}

func GetBroaderMimeType(mime string) string {
	switch {
	case strings.HasPrefix(mime, "text/"):
		return "text/plain"
	default:
		return ""
	}
}
