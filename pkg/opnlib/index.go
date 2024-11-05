package opnlib

import (
	"encoding/json"
	"fmt"
	"github.com/MatthiasKunnen/xdg/desktop"
	"github.com/MatthiasKunnen/xdg/mimeapps"
	"os"
	"path"
	"time"
)

type Index struct {
	Version     int
	GeneratedOn time.Time

	// Associations is a map of Key=MIME type, Value=List of desktop IDs.
	// It can be used to look up all the desktop IDs that support opening a certain MIME type.
	Associations mimeapps.Associations

	// DesktopIdToPaths maps a desktop ID to the desktop files in the file system.
	DesktopIdToPaths desktop.IdPathMap

	// isNewlyGenerated is true when the index is not loaded from cache.
	isNewlyGenerated bool
}

// LoadIndex loads the index.
func LoadIndex(filename string) (*Index, error) {
	content, err := os.ReadFile(filename)

	if err != nil {
		return nil, fmt.Errorf("error loading index from '%s': %w", filename, err)
	}

	var db Index
	err = json.Unmarshal(content, &db)
	if err != nil {
		return nil, fmt.Errorf("parsing error loading index from '%s': %w", filename, err)
	}

	return &db, nil
}

// SaveIndex saves the index into the given file.
func (index *Index) SaveIndex(filename string) error {
	err := os.MkdirAll(path.Dir(filename), 0750)
	if err != nil {
		return err
	}

	var file *os.File
	file, err = os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file at %s: %w", filename, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(index)
	if err != nil {
		return fmt.Errorf("error saving JSON data at %s: %w", filename, err)
	}

	return nil
}

// GenerateIndex generates the index used to look up MIME type/Application associations and the
// paths of desktop files.
func GenerateIndex() (*Index, error) {
	locations := desktop.GetDesktopFileLocations()
	idPathMap, err := desktop.GetDesktopFiles(locations)
	if err != nil {
		return nil, fmt.Errorf("error getting desktop files: %w", err)
	}

	currentDesktop := os.Getenv("XDG_CURRENT_DESKTOP")
	lists := mimeapps.GetLists(currentDesktop)
	associations := mimeapps.GetPreferredApplications(lists, idPathMap)

	return &Index{
		Version:          1,
		GeneratedOn:      time.Now(),
		Associations:     associations,
		DesktopIdToPaths: idPathMap,
		isNewlyGenerated: true,
	}, nil
}
