package main

import (
	"github.com/MatthiasKunnen/opn/cmd"
)

func main() {
	cmd.Execute()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//mimeappsLists := mimeapps.GetLists(os.Getenv("XDG_CURRENT_DESKTOP"))
	//fmt.Printf("mimeapps: %#v\n", mimeappsLists)
	//
	//file, err := os.Open("/home/matthias/.config/mimeapps.list")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer file.Close()
	//
	//locations := desktop.GetDesktopFileLocations()
	//idPathMap, err := desktop.GetDesktopFiles(locations)
	//
	//_, path, err := basedir.CreateSystemDataFile("opn/db.json")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("created system data file: %s\n", path)
	//
	//lists, err := mimeapps.ParseMimeapps(file)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("mimeapps parsed: %v\n", lists)
}
