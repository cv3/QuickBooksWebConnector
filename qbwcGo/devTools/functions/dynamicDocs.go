package devTools

import (
	"bytes"
	"fmt"
	"html/template"
	"os"

	"github.com/amazingfly/cv3go"
)

//BuildDynamicDocumentation will read in all the mapping files and create an html document
func BuildDynamicDocumentation(mapPath, tPath string) {
	var allMaps = make(map[string]interface{}, 0)
	var templateBuf = bytes.Buffer{}

	ScanAllFilesInDir(mapPath, func(f os.FileInfo) {
		if !f.IsDir() {
			if f.Name() != "stateTaxMapping.json" {
				mapObjects := ReadRobustFieldMapping(mapPath + f.Name())
				allMaps[f.Name()] = mapObjects
				fmt.Println("appending: ", f.Name())
			}
		}
	})
	fmt.Println("===================")
	fmt.Println(allMaps)
	fmt.Println("===================")
	LoadTemplate(&tPath, &allMaps, &templateBuf)
	fmt.Println("template done")
	cv3go.PrintToFile(templateBuf.Bytes(), "testTables.html")
}

//LoadTemplate will take a path a data struct and a bytes.Buffer tehn load the template
func LoadTemplate(tPath *string, ctx interface{}, requestBody *bytes.Buffer) {
	t, err := template.ParseFiles(*tPath)
	if err != nil {
		fmt.Println(err)
	} // Populate requestBody with the executed template and context

	err = t.Execute(requestBody, ctx)
	if err != nil {
		fmt.Println(err)
	}
}
