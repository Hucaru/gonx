package gonx

import (
	"testing"
)

func TestFile(t *testing.T) {
	_, _, _, _, err := Parse("testData/Data_NoBlobs_PKG4.nx")

	if err != nil {
		panic(err)
	}
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _, _, _, err := Parse("testData/Data_NoBlobs_PKG4.nx")

		if err != nil {
			panic(err)
		}
	}
}
