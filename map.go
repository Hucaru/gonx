package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Map data from nx
type Map struct {
}

// ExtractMaps from parsed nx
func ExtractMaps(nodes []Node, textLookup []string) map[int32]Map {
	maps := make(map[int32]Map)

	searches := []string{"/Map/Map/Map0", "/Map/Map/Map1", "/Map/Map/Map2", "/Map/Map/Map9"}

	for _, search := range searches {
		valid := searchNode(search, nodes, textLookup, func(node *Node) {
			for i := uint32(0); i < uint32(node.ChildCount); i++ {
				mapNode := nodes[node.ChildID+i]
				name := textLookup[mapNode.NameID]

				subSearches := []string{"info", "life", "portal", "reactor"}
				var mapItem Map

				for _, s := range subSearches {
					newSearch := search + "/" + name + "/" + s
					// Reactor is not always present so ignore the valid return from search
					searchNode(newSearch, nodes, textLookup, func(node *Node) {
						addToMap(&mapItem, node, nodes, textLookup)
					})
				}

				name = strings.TrimSuffix(name, filepath.Ext(name))
				mapID, err := strconv.Atoi(name)

				if err != nil {
					log.Println(err)
					continue
				}

				maps[int32(mapID)] = mapItem
			}
		})

		if !valid {
			log.Println("Invalid node search:", search)
		}

	}

	return maps
}

func addToMap(mapItem *Map, node *Node, nodes []Node, textLookup []string) {
	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		option := nodes[node.ChildID+i]
		optionName := textLookup[option.NameID]

		switch optionName {
		default:
			// log.Println("Unsupported NX map option:", optionName, "->", option.Data)
		}
	}
}
