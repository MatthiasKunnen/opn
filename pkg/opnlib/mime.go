package opnlib

import (
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"
)

// GetFileMime returns the mime type of the given path.
// It uses programs installed
func GetFileMime(path string) (string, error) {
	xdgCmd := exec.Command("xdg-mime", "query", "filetype", path)
	output, err := xdgCmd.Output()

	var exitError *exec.ExitError

	switch {
	case errors.Is(err, exec.ErrNotFound):
		// xdg-mime not found in PATH. Fall back to file.
	case errors.As(err, &exitError):
		switch exitError.ExitCode() {
		case 2:
			// See Exit Codes in XDG-MIME(1)
			return "", fs.ErrNotExist
		default:
			return "", fmt.Errorf("xdg-mime exited with %d: %s", exitError.ExitCode(), string(output))
		}
	default:
		return strings.TrimSpace(string(output)), nil
	}

	fileCmd := exec.Command(
		"/usr/bin/file",
		"-E",            // Exit with non-zero exit code on filesystem errors
		"--brief",       // Do not prepend filenames to output
		"--dereference", // Follow symlinks
		"--mime-type",   // Print MIME type only
		path,
	)
	output, err = fileCmd.Output()

	switch {
	case errors.Is(err, exec.ErrNotFound), errors.Is(err, fs.ErrNotExist):
		return "", fmt.Errorf("failed to determine MIME type, no programs to determine" +
			" MIME type are installed. Either xdg-mime (xdg-utils) or, file, is required")
	case errors.As(err, &exitError):
		return "", fmt.Errorf("file exited with %d: %s", exitError.ExitCode(), string(output))
	case err != nil:
		return "", err
	default:
		return strings.TrimSpace(string(output)), nil
	}
}
