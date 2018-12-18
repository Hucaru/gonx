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
	HP              int16
	Limit, Interval int16
	MobID           []int32
}

// ExtractSkills from parsed nx
func ExtractSkills(nodes []Node, textLookup []string) (map[int32][]PlayerSkill, map[int32][]MobSkill) {
	playerSkills := make(map[int32][]PlayerSkill)
	mobSkills := make(map[int32][]MobSkill)

	search := "/Skill"

	valid := searchNode(search, nodes, textLookup, func(node *Node) {
		for i := uint32(0); i < uint32(node.ChildCount); i++ {
			skillSectionNode := nodes[node.ChildID+i]
			name := textLookup[skillSectionNode.NameID]

			if _, err := strconv.Atoi(strings.TrimSuffix(name, filepath.Ext(name))); err != nil {
				mobSkillSearch := search + "/" + name
				skillIDs := []string{}

				valid := searchNode(mobSkillSearch, nodes, textLookup, func(node *Node) {
					for j := uint32(0); j < uint32(node.ChildCount); j++ {
						skillNode := nodes[node.ChildID+j]
						skillIDs = append(skillIDs, textLookup[skillNode.NameID])
					}
				})

				for _, s := range skillIDs {
					valid = searchNode(mobSkillSearch+"/"+s+"/level", nodes, textLookup, func(node *Node) {
						skillID, err := strconv.Atoi(s)

						if err != nil {
							return
						}

						mobSkills[int32(skillID)] = make([]MobSkill, node.ChildCount)

						for j := uint32(0); j < uint32(node.ChildCount); j++ {
							skillNode := nodes[node.ChildID+j]
							skillLevel := textLookup[skillNode.NameID]
							level, err := strconv.Atoi(skillLevel)

							if err == nil {
								mobSkills[int32(skillID)][level-1] = getMobSkill(&skillNode, nodes, textLookup)
							}
						}
					})
				}

				if !valid {
					log.Println("Invalid node search:", mobSkillSearch)
				}
			} else {
				playerSkillSearch := search + "/" + name + "/skill"
				skillIDs := []string{}

				valid := searchNode(playerSkillSearch, nodes, textLookup, func(node *Node) {
					for j := uint32(0); j < uint32(node.ChildCount); j++ {
						skillNode := nodes[node.ChildID+j]
						skillIDs = append(skillIDs, textLookup[skillNode.NameID])
					}
				})

				for _, s := range skillIDs {
					valid = searchNode(playerSkillSearch+"/"+s+"/level", nodes, textLookup, func(node *Node) {
						skillID, err := strconv.Atoi(s)

						if err != nil {
							return
						}

						playerSkills[int32(skillID)] = make([]PlayerSkill, node.ChildCount)

						for j := uint32(0); j < uint32(node.ChildCount); j++ {
							skillNode := nodes[node.ChildID+j]
							skillLevel := textLookup[skillNode.NameID]
							level, err := strconv.Atoi(skillLevel)

							if err == nil {
								playerSkills[int32(skillID)][level-1] = getPlayerSkill(&skillNode, nodes, textLookup)
							}
						}
					})
				}

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

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		option := nodes[node.ChildID+i]
		optionName := textLookup[option.NameID]

		switch optionName {
		case "mad":
		case "mdd":
		case "pad":
		case "pdd":
		case "hp":
		case "mp":
		case "hpCon":
		case "mpCon":
		case "bulletConsume":
		case "moneyCon":
		case "itemCon":
		case "itemConNo":
		case "mastery":
		case "time":
		case "eva":
		case "acc":
		case "jump":
		case "speed":
		case "range":
		case "mobCount":
		case "attackCount":
		case "damage":

		// unsure what these are for
		case "fixdamage":
		case "rb":
		case "hs":
		case "lt":
		case "x":
		case "y":
		case "z":
		case "prop":
		case "ball":
		case "hit":
		case "bulletCount":
		case "action":
		case "58": //?
		default:
			log.Println("Unsupported NX player skill option:", optionName, "->", option.Data)
		}
	}

	return skill
}

func getMobSkill(node *Node, nodes []Node, textLookup []string) MobSkill {
	skill := MobSkill{}

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		option := nodes[node.ChildID+i]
		optionName := textLookup[option.NameID]

		switch optionName {
		case "hp":
		case "interval":
		case "limit":
		case "summonEffect":
		case "time":
		case "mpCon":

		// ?
		case "0":
		case "1":
		case "2":
		case "3":
		case "4":
		case "5":

		// not sure what these are used for
		case "lt":
		case "rb":
		case "effect":
		case "x":
		case "y":
		case "tile":
		case "prop":
		case "affected":
		case "mob":
		case "mob0":
		default:
			log.Println("Unsupported NX mob skill option:", optionName, "->", option.Data)
		}
	}

	return skill
}
