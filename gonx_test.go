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

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		if *nxFile == "" {
			log.Fatal("No NX file specified")
		}
		_, _, err := Parse(*nxFile)

		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkExtract(b *testing.B) {
	if *nxFile == "" {
		log.Fatal("No NX file specified")
	}
	nodes, textLookup, err := Parse(*nxFile)

	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		ExtractMobs(nodes, textLookup)
		ExtractMaps(nodes, textLookup)
		ExtractItems(nodes, textLookup)
		ExtractSkills(nodes, textLookup)
	}
}
