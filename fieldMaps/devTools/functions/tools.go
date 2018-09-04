package devTools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/TeamFairmont/QuickBooksWebConnector/qbwcGo"
	"github.com/amazingfly/cv3go"
)

//ScanForCV3Fields will look through the robust mapping files and let me know the fiel and field that has a populated cv3Field
func ScanForCV3Fields(mapFilesDir string) {
	files, err := ioutil.ReadDir(mapFilesDir)
	if err != nil {
		fmt.Println("Error reading directory of map files in BuildNewMappingFiles: ", err)
	} else { //no error
		for _, f := range files {
			if !f.IsDir() {
				fmt.Println(mapFilesDir + f.Name())
				fBytes, err := ioutil.ReadFile(mapFilesDir + f.Name())
				var mapObjects = make(map[string]qbwcGo.MappingObject, 0)
				err = json.Unmarshal(fBytes, &mapObjects)
				if err != nil {
					fmt.Println("error unmarshalling map file: ", err)
				}
				fmt.Println("=====================")
				for key, mObj := range mapObjects {
					for _, mData := range mObj {
						if mData.Data != "" {
							fmt.Println(key + ": " + mData.Data + " in " + f.Name())
						}
					}
				}
				fmt.Println("=====================")
			}
		}
	}
}

//BuildNewMappingFiles builds the new syle map files from the old style mapping files from a passed in directory
func BuildNewMappingFiles(mapFilesDir string) {
	files, err := ioutil.ReadDir(mapFilesDir)
	if err != nil {
		fmt.Println("Error reading directory of map files in BuildNewMappingFiles: ", err)
	} else { //no error

		for _, f := range files {

			if !f.IsDir() {
				//files = append(files[:i], files[i+1:]...)
				simpleMap := ReadFieldMapping(mapFilesDir + f.Name())
				robustMap := BuildRobustMapping(simpleMap)

				mapBytes, err := json.MarshalIndent(&robustMap, "", "\t")
				if err != nil {
					fmt.Println("Error marshalling json from robustMap in BuildNewMappingFiles: ", err)
				} else { //no error
					cv3go.PrintToFile(mapBytes, mapFilesDir+f.Name())
				}
			}
		}
	}
}

//BuildRobustMapping will take a map[string]string and build a map[string]MappingObject.  Origionally for converting old mapping files to the new, more robust mapping objects
func BuildRobustMapping(simpleMap map[string]string) map[string]qbwcGo.MappingObject {
	var robustMap = make(map[string]qbwcGo.MappingObject, len(simpleMap))

	for qbField, cv3Field := range simpleMap {
		var mapObject = qbwcGo.MappingObject{}
		var mapData = qbwcGo.MapData{
			Data:        cv3Field,
			MappedField: true,
		}
		mapObject = append(mapObject, mapData)
		robustMap[qbField] = mapObject
	}
	return robustMap
}

//ReadFieldMapping reads the json field map file and returns a map
func ReadFieldMapping(mapPath string) map[string]string {
	var fieldMap map[string]string
	mapFile, err := ioutil.ReadFile(mapPath)
	if err != nil {
		fmt.Println("error loading map file: ", err)
	}
	err = json.Unmarshal(mapFile, &fieldMap)
	if err != nil {
		fmt.Println("error unmarshalling map file: ", err)
	}
	return fieldMap
}

//ReadRobustFieldMapping reads the json field map file and returns a map
func ReadRobustFieldMapping(mapPath string) map[string]qbwcGo.MappingObject {
	var fieldMap map[string]qbwcGo.MappingObject
	mapFile, err := ioutil.ReadFile(mapPath)
	if err != nil {
		fmt.Println("error loading map file: ", err)
	}
	err = json.Unmarshal(mapFile, &fieldMap)
	if err != nil {
		fmt.Println("error unmarshalling map file: ", err)
	}
	return fieldMap
}

//ScanAllFilesInDir will look through all files in a given directory and execute the passed in function
func ScanAllFilesInDir(mapFilesDir string, job func(os.FileInfo)) {
	files, err := ioutil.ReadDir(mapFilesDir)
	if err != nil {
		fmt.Println("Error reading directory of map files in BuildNewMappingFiles: ", err)
	} else { //no error

		for _, f := range files {
			job(f)
		}
	}
}
