package gonx

// Reactor data from nx
type Reactor struct {
}

// ExtractReactors from parsed nx
func ExtractReactors(nodes []Node, textLookup []string) map[int32]Reactor {
	return make(map[int32]Reactor)
}
