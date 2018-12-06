package gonx

import (
	"flag"
	"testing"
)

var nxFile = flag.String("nxFile", "", "Path to nx file to use for testing")

func TestFile(t *testing.T) {
	nodes, textLookup := Parse(*nxFile)

	ExtractMobs(nodes, textLookup)
	ExtractMaps(nodes, textLookup)
	ExtractItems(nodes, textLookup)
}
