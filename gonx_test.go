package gonx

import (
	"flag"
	"log"
	"testing"
)

var nxFile = flag.String("nxFile", "", "Path to nx file to use for testing")

func TestFile(t *testing.T) {
	if *nxFile == "" {
		log.Fatal("No NX file specified")
	}
	nodes, textLookup, err := Parse(*nxFile)

	if err != nil {
		panic(err)
	}

	ExtractMobs(nodes, textLookup)
	ExtractMaps(nodes, textLookup)
	ExtractItems(nodes, textLookup)
	ExtractSkills(nodes, textLookup)
}
