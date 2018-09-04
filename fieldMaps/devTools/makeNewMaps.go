/*
BACKUP YOUR FILES FIRST
This is a tool to convert the old style maps into a more robust mapping style.
It looks in the given directory and converts json files of type map[string]string into json files of type map[string]qbwcGo.MappingObject

Command line arguement -dir accepts a directory string, the default directory is "../../fieldMaps/"

*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	devTools "github.com/TeamFairmont/QuickBooksWebConnector/qbwcGo/devTools/functions"
)

func main() {
	//command line flags
	var dir = flag.String("dir", "../../fieldMaps/", "directory of the json map[string]string files to be converted")
	flag.Parse()
	_ = dir

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("This will overwrite the map files in the directory: ", *dir)
	fmt.Println("continue?(y/n)")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading user input: ", err)
	}
	if strings.Contains(strings.ToLower(input), "y") {
		devTools.BuildNewMappingFiles(*dir)
		fmt.Println("maps converted in dir: ", *dir)
	} else {
		fmt.Println("No files changed")
	}
	fmt.Println("======================")
	devTools.ScanForCV3Fields(*dir)

	//devTools.BuildDynamicDocumentation()
}
