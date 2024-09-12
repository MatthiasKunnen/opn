package main

import (
	"fmt"
	"github.com/MatthiasKunnen/xdg/basedir"
	"github.com/MatthiasKunnen/xdg/desktop"
	"github.com/MatthiasKunnen/xdg/mimeapps"
	"log"
	"os"
)

func main() {
	mimeappsLists := mimeapps.GetLists(os.Getenv("XDG_CURRENT_DESKTOP"))
	fmt.Printf("mimeapps: %#v\n", mimeappsLists)

	file, err := os.Open("/home/matthias/.config/mimeapps.list")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	idPathMap, err := desktop.GetDesktopFiles()

	_, path, err := basedir.CreateSystemDataFile("opn/db.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("created system data file: %s\n", path)

	lists, err := mimeapps.ParseMimeapps(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("mimeapps parsed: %v\n", lists)

}

// removeDuplicates removes duplicates entries from a slice and returns the slice.
// Order is preserved and the first occurrence of every entry is preserved.
func removeDuplicates[T comparable](input []T) []T {
	seen := make(map[T]bool, len(input))
	list := make([]T, 0, len(input))

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			list = append(list, item)
		}
	}

	return list
}
