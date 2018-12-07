# gonx
Parse NX files and extract data into useable structures

## Benchmarks:

To run the benchmarks type go test -nxFile ../Data.nx -run=XXX -bench=.

```
goos: windows
goarch: amd64
pkg: github.com/Hucaru/gonx
BenchmarkParse-4               5         264326940 ns/op
BenchmarkExtract-4            50          24658304 ns/op
PASS
ok      github.com/Hucaru/gonx  6.375s
```