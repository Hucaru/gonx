# gonx
Parse NX files and extract data into useable structures (only works on version of the game with Data.wz).

## Usage

```golang
package main

import "github.com/Hucaru/gonx"

func main() {
    nodes, textLookup, bitmaps, audio, err := gonx.Parse(fname)

    if err != nil {
        panic(err)
    }

    found := gonx.FindNode("element/childEment/grandchildElement", func(node *gonx.Node) {
        // if image
        img, x, y, err := node.GetBitmap(bitmaps)

        // if audio
        sound, err := node.GetAudio(audio)

        // node.Type - 1 = int64, 2 = float64
        data = node.Data
    })

    if !found {
        panic(fmt.Errorf("Not a valid node"))
    }
}
```

## Benchmarks
To run the benchmark tests type `go test -run=XXX -bench=.`

i5-3570k @ 3.40GHz
```
goos: windows
goarch: amd64
pkg: github.com/Hucaru/gonx
BenchmarkParse-4              20          86080150 ns/op
BenchmarkExtract-4            30          33965990 ns/op
PASS
ok      github.com/Hucaru/gonx  4.463s
```

i5 @ 2.40GHz
```
goos: darwin
goarch: amd64
pkg: github.com/Hucaru/gonx
BenchmarkParse-4   	      10	 113918050 ns/op
PASS
ok  	github.com/Hucaru/gonx	2.040s
```