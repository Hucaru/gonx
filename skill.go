package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// PlayerSkill data from nx
type PlayerSkill struct {
	MaxLevel byte
	Mastery  []int16
}

// MobSkill data from nx
type MobSkill struct {
}

// ExtractSkills from parsed nx
func ExtractSkills(nodes []Node, textLookup []string) (map[int32][]PlayerSkill, map[int32][]MobSkill) {
	playerSkills := make(map[int32][]PlayerSkill)
	mobSkills := make(map[int32][]MobSkill)

	search := "/Skill"

	valid := searchNode(search, nodes, textLookup, func(node *Node) {
		for i := uint32(0); i < uint32(node.ChildCount); i++ {
			mapNode := nodes[node.ChildID+i]
			name := textLookup[mapNode.NameID]

			if _, err := strconv.Atoi(strings.TrimSuffix(name, filepath.Ext(name))); err != nil {
				mobSkillSearch := search + "/" + name
				valid := searchNode(mobSkillSearch, nodes, textLookup, func(node *Node) {
					// loop over all skills
				})

				if !valid {
					log.Println("Invalid node search:", mobSkillSearch)
				}
			} else {
				playerSkillSearch := search + "/" + name
				valid := searchNode(playerSkillSearch, nodes, textLookup, func(node *Node) {
					// loop over all skills
				})

				if !valid {
					log.Println("Invalid node search:", playerSkillSearch)
				}
			}

		}
	})

	if !valid {
		log.Println("Invalid node search:", search)
	}

	return playerSkills, mobSkills
}

func getPlayerSkill(node *Node, nodes []Node, textLookup []string) PlayerSkill {
	skill := PlayerSkill{}

	return skill
}

func getMobSkill(node *Node, nodes []Node, textLookup []string) MobSkill {
	skill := MobSkill{}

	return skill
}
