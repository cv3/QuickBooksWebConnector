package main

import (
	"flag"

	devTools "github.com/TeamFairmont/QuickBooksWebConnector/qbwcGo/devTools/functions"
)

func main() {
	//command line flags
	var dir = flag.String("dir", "../../fieldMaps/", "directory of the QBWC mapping files")
	var tPath = flag.String("tPath", "../../templates/dynamicDocumentationHTML.t", "path to the dynamic documentation template")
	flag.Parse()
	_ = dir

	devTools.BuildDynamicDocumentation(*dir, *tPath)
}
