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

// GetBroaderMimeType returns a mime type that matches the given type more broadly.
// If no broader mime type is known, an empty string is returned.
// E.g.
// - application/javascript returns text/plain.
// - text/csv returns text/plain.
func GetBroaderMimeType(mime string) string {
	// @TODO we might want to make this configurable.
	// In theory, users could add all the mime→desktop id mappings manually but this would be
	// cumbersome. Providing sensible defaults is probably the right choice but might not be perfect
	// for everyone. Perhaps there could be a config file that overrides the defaults here.
	switch mime {
	case
		"application/javascript",
		"application/json",
		"application/ld+json",
		"application/xml",
		"application/yaml",
		"image/svg+xml":
		return "text/plain"
	case "text/plain":
		return ""
	default:
		if strings.HasPrefix(mime, "text/") {
			return "text/plain"
		}

		return ""
	}
}
