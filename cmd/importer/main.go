package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	alabamaFilePath := flag.String("alabama", "./gpx/ACS.gpx", "File path to the ACS gpx file to import")
	dsn := flag.String("dsn", "", "REQUIRED data source name for mysql driver")
	flag.Parse()

	if dsn == nil || *dsn == "" {
		fmt.Println("dsn is invalid")
		os.Exit(1)
	}

	alabamaReader, err := os.Open(*alabamaFilePath)
	if err != nil {
		fmt.Println("unable to open ACS file")
		os.Exit(1)
	}

	// Everything above is not unit-tested
	if err = parseAndImportCaves(alabamaReader, *dsn); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Execution completed successfully")
	}
}
