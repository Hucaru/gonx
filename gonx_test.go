package gonx

import (
	"flag"
	"testing"
)

var nxFile = flag.String("nxFile", "", "Path to nx file to use for testing")

func TestFile(t *testing.T) {
	nodes, textLookup := Parse(*nxFile)

	_ = ExtractMobs(nodes, textLookup)
	_ = ExtractMaps(nodes, textLookup)
	_ = ExtractItems(nodes, textLookup)
}
