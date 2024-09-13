package opnlib

import (
	"encoding/json"
	"fmt"
	"github.com/MatthiasKunnen/xdg/basedir"
	"github.com/MatthiasKunnen/xdg/desktop"
	"github.com/MatthiasKunnen/xdg/mimeapps"
	"os"
)

type Db struct {
	Version          int
	Associations     mimeapps.Associations
	DesktopIdToPaths desktop.IdPathMap
}

const dbLocation = "opn/db.json"

var (
	ErrDbNotFound = fmt.Errorf("%s not found", dbLocation)
)

// LoadDb loads the database at the given file or the first one that can be found in the
func LoadDb(path string) (*Db, error) {
	if path == "" {
		foundPath, err := basedir.FindDataFile(dbLocation)
		if err != nil {
			return nil, fmt.Errorf("error locating db.json: %w", err)
		} else if foundPath == "" {
			return nil, ErrDbNotFound
		}
		path = foundPath
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading '%s': %w", path, err)
	}

	var db Db
	err = json.Unmarshal(content, &db)
	if err != nil {
		return nil, fmt.Errorf("error converting JSON read from '%s': %w", path, err)
	}

	return &db, nil
}

func MustLoadDb(path string) *Db {
	// 1.Attempt to load db
	// 3. if success, return
	// 2. if fail, generate and persist, return generated db
}
