package opnlib

import (
	"encoding/json"
	"fmt"
	"github.com/MatthiasKunnen/xdg/basedir"
	"github.com/MatthiasKunnen/xdg/desktop"
	"github.com/MatthiasKunnen/xdg/mimeapps"
	"log"
	"os"
	"path/filepath"
)

type Index struct {
	Version          int
	Associations     mimeapps.Associations
	DesktopIdToPaths desktop.IdPathMap
}

const CacheLocation = "opn/db.json"

var (
	ErrIndexNotFound = fmt.Errorf("%s not found", CacheLocation)
)

// LoadIndex loads the index at the given filename.
//   - If the filename is an absolute path, it is loaded.
//   - If the filename is a relative path, it is searched for in the config directories and loaded.
//     If it is not found in the config directories, an [ErrIndexNotFound] is returned.
//   - If the filename is empty, it is defaulted to [CacheLocation] and the relative path behavior is
//     used.
//
// Other file system error can occur and will be returned.
func LoadIndex(filename string) (*Index, error) {
	if filename == "" {
		filename = CacheLocation
	}

	var path string
	var err error
	if !filepath.IsAbs(filename) {
		path, err = basedir.FindDataFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error locating %s: %w", filename, err)
		} else if path == "" {
			return nil, ErrIndexNotFound
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading '%s': %w", path, err)
	}

	var db Index
	err = json.Unmarshal(content, &db)
	if err != nil {
		return nil, fmt.Errorf("error converting JSON read from '%s': %w", path, err)
	}

	return &db, nil
}

// MustLoadIndex loads the index and returns it.
//
// When skipCache is true, loading of the index is skipped. Instead, it is generated and returned.
//
// When skipCache is false, loading of the index is attempted at the filename's location:
//   - If the filename is an absolute path, it is loaded.
//   - If the filename is a relative path, it is searched for in the config directories and loaded.
//     If it is not found in the config directories, it is generated and saved.
//   - If the filename is empty, it is defaulted to [CacheLocation] and the relative path behavior is
//     used.
//
// If loading fails, the index is generated and saved.
//
// This function panics when:
//   - skipCache is true and the index could not be generated.
//   - skipCache is false and the index could not be generated and saved.
func MustLoadIndex(skipCache bool, filename string) *Index {
	if skipCache {
		db, err := GenerateIndex()
		if err != nil {
			log.Fatalf("failed to generate index: %v", err)
		}

		return db
	}

	if filename == "" {
		filename = CacheLocation
	}

	var path string
	var err error
	if !filepath.IsAbs(filename) {
		path, err = basedir.FindDataFile(filename)
	}

	if err == nil {
		content, err := os.ReadFile(path)
		if err == nil {
			var db Index
			err = json.Unmarshal(content, &db)
			if err == nil {
				return &db
			}
		}
	}

	index, err := GenerateIndex()
	if err != nil {
		log.Fatalf("failed to load existing index.json and failed to generate it: %v", err)
	}

	err = index.SaveIndex(filename)
	if err != nil {
		log.Fatalf(
			"failed to load existing index.json and failed to save a newly generated one: %v",
			err,
		)
	}

	return index
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
		Associations:     associations,
		DesktopIdToPaths: idPathMap,
	}, nil
}

// SaveIndex saves the index into a cache file.
// If filename is an absolute path, the index is saved to it.
// If filename is not absolute, the index is saved into the first writable config directory.
// If filename is empty, the non-absolute behavior occurs defaulting to [CacheLocation].
func (index *Index) SaveIndex(filename string) error {
	if filename == "" {
		filename = CacheLocation
	}

	var file *os.File
	var err error
	var path string
	if filepath.IsAbs(filename) {
		path = filename
		file, err = os.Create(filename)
	} else {
		file, path, err = basedir.CreateDataFile(filename)
	}
	if err != nil {
		return fmt.Errorf("error creating file at %s: %w", path, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(index)
	if err != nil {
		return fmt.Errorf("error saving JSON data at %s: %w", path, err)
	}

	return nil
}
