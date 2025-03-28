package opnlib

import (
	"errors"
	"fmt"
	"github.com/MatthiasKunnen/xdg/basedir"
	"github.com/MatthiasKunnen/xdg/sharedmimeinfo"
	"path"
	"time"
)

type Opn struct {
	// CacheFilePath is either an absolute path to the cache file or a path relative to the cache
	// dir. Leave empty to use the default.
	CacheFilePath string

	// index holds the lookup maps for associations and desktop IDs.
	index *Index

	// SkipCache determines whether loading the cache is skipped.
	SkipCache bool

	mimeSubclassInfo *sharedmimeinfo.Subclass
}

var (
	FailedToSaveCache = errors.New("failed to save cache")
)

// Load loads the cache and generates it if necessary.
func (opn *Opn) Load() error {
	filename := opn.getCachePath()

	if opn.SkipCache {
		index, err := GenerateIndex()
		if err != nil {
			return fmt.Errorf("failed to generate index: %w", err)
		}

		opn.index = index
		return nil
	}

	index, err := LoadIndex(filename)
	if err == nil && index.GeneratedOn.Add(24*time.Hour).After(time.Now()) {
		opn.index = index
	} else {
		index, err = GenerateIndex()
		if err != nil {
			return fmt.Errorf("failed to generate index: %w", err)
		}

		opn.index = index
	}

	opn.mimeSubclassInfo, err = sharedmimeinfo.LoadFromOs()
	if err != nil {
		return fmt.Errorf("failed to load subclass info: %w", err)
	}

	return nil
}

// LoadAndSave attempts to load the cache and saves it if necessary.
func (opn *Opn) LoadAndSave() error {
	err := opn.Load()
	if err != nil {
		return err
	}

	if opn.index.isNewlyGenerated {
		err = opn.SaveIndex()
		if err != nil {
			return err
		}
	}

	return nil
}

type MimeDesktopIds struct {
	Mime       string
	DesktopIds []string
}

// GetDesktopIdsForBroadMime returns all desktop IDs for a given mime type and all its subtypes.
// E.g. text/html will also return results for text/plain.
// The results are in order of higher priority to lower priority.
func (opn *Opn) GetDesktopIdsForBroadMime(mimeType string) []MimeDesktopIds {
	result := []MimeDesktopIds{
		{
			Mime:       mimeType,
			DesktopIds: opn.GetDesktopIdsForMime(mimeType),
		},
	}

	broaderMime := opn.mimeSubclassInfo.BroaderDfs(mimeType)
	for _, mime := range broaderMime {
		result = append(result, MimeDesktopIds{
			Mime:       mime,
			DesktopIds: opn.GetDesktopIdsForMime(mime),
		})
	}

	return result
}

// GetDesktopIdsForMime returns all desktop IDs for a given mime type.
// The results are in order of higher priority to lower priority.
func (opn *Opn) GetDesktopIdsForMime(mimeType string) []string {
	associations := opn.index.Associations[mimeType]
	associationsCopy := make([]string, len(associations))
	copy(associationsCopy, associations)

	return associationsCopy
}

func (opn *Opn) GetDesktopFileLocations(desktopId string) []string {
	return opn.index.DesktopIdToPaths[desktopId]
}

func (opn *Opn) getCachePath() string {
	var filename = opn.CacheFilePath

	if filename == "" {
		filename = GetDefaultCachePath()
	}

	if !path.IsAbs(filename) {
		return path.Join(basedir.CacheHome, filename)
	}

	return filename
}

func (opn *Opn) SaveIndex() error {
	err := opn.index.SaveIndex(opn.getCachePath())
	if err != nil {
		return fmt.Errorf("%w: %w", FailedToSaveCache, err)
	}

	return nil
}

func GetDefaultCachePath() string {
	return path.Join(basedir.CacheHome, "opn/db.json")
}
