# gonx
Parse NX files and extract data into useable structures (only works on version of the game with Data.wz).

## Usage

```golang
package main

import "github.com/Hucaru/gonx"

func main() {
    nodes, textLookup, err := gonx.Parse(fname)

    if err != nil {
        panic(err)
    }

    mobs := gonx.ExtractMobs(nodes, textLookup)
    maps := gonx.ExtractMaps(nodes, textLookup)
    items := gonx.ExtractItems(nodes, textLookup)
    playerSkills, mobSkills := gonx.ExtractSkills(nodes, textLookup)
}
```

## Benchmarks
On i5-3570k with a clock speed of 3.40GHz it takes < 0.4 seconds to parse and extract all the data from an a v28 nx file without media.

To run the benchmark tests type `go test -nxFile ../Data.nx -run=XXX -bench=.`
```
goos: windows
goarch: amd64
pkg: github.com/Hucaru/gonx
BenchmarkParse-4              20          86129320 ns/op
BenchmarkExtract-4            30          33999273 ns/op
PASS
ok      github.com/Hucaru/gonx  4.524s
```