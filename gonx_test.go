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
	nodes, textLookup := Parse(*nxFile)

	ExtractMobs(nodes, textLookup)
	ExtractMaps(nodes, textLookup)
	ExtractItems(nodes, textLookup)
	ExtractSkills(nodes, textLookup)
}
