package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Mob data from nx
type Mob struct {
}

// ExtractMobs from parsed nx
func ExtractMobs(nodes []Node, textLookup []string) map[int32]Mob {
	mobs := make(map[int32]Mob)

	search := "/Mob"
	valid := searchNode(search, nodes, textLookup, func(node *Node) {
		for i := uint32(0); i < uint32(node.ChildCount); i++ {
			mobNode := nodes[node.ChildID+i]
			name := textLookup[mobNode.NameID]
			subSearch := search + "/" + name + "/info"

			var mob Mob

			valid := searchNode(subSearch, nodes, textLookup, func(node *Node) {
				mob = getMob(node, nodes, textLookup)
			})

			if !valid {
				log.Println("Invalid node search:", subSearch)
			}

			name = strings.TrimSuffix(name, filepath.Ext(name))
			itemID, err := strconv.Atoi(name)

			if err != nil {
				log.Println(err)
				continue
			}

			mobs[int32(itemID)] = mob
		}
	})

	if !valid {
		log.Println("Invalid node search:", search)
	}

	return mobs
}

func getMob(node *Node, nodes []Node, textLookup []string) Mob {
	mob := Mob{}

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		option := nodes[node.ChildID+i]
		optionName := textLookup[option.NameID]

		switch optionName {
		default:
			log.Println("Unsupported NX mob option:", optionName, "->", option.Data)
		}
	}
}
