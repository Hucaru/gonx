package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// PlayerSkill data from nx
type PlayerSkill struct {
	Mastery               int64
	Mad, Mdd, Pad, Pdd    int64
	Hp, Mp, HpCon, MpCon  int64
	BulletConsume         int64
	MoneyConsume          int64
	ItemCon               int64
	ItemConNo             int64
	Time                  int64
	Eva, Acc, Jump, Speed int64
	Range                 int64
	MobCount              int64
	AttackCount           int64
	Damage                int64
	Fixdamage             int64
	Rb, Lt                Vector
	Hs                    string
	X, Y, Z               int64
	Prop                  int64
	BulletCount           int64
	Action                string
}

// MobSkill data from nx
type MobSkill struct {
	HP              int64
	Limit, Interval int64
	MobID           []int64
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
			skill.Mad = dataToInt64(option.Data)
		case "mdd":
			skill.Mdd = dataToInt64(option.Data)
		case "pad":
			skill.Pad = dataToInt64(option.Data)
		case "pdd":
			skill.Pdd = dataToInt64(option.Data)
		case "hp":
			skill.Hp = dataToInt64(option.Data)
		case "mp":
			skill.Mp = dataToInt64(option.Data)
		case "hpCon":
			skill.HpCon = dataToInt64(option.Data)
		case "mpCon":
			skill.MpCon = dataToInt64(option.Data)
		case "bulletConsume":
			skill.BulletConsume = dataToInt64(option.Data)
		case "moneyCon":
			skill.MoneyConsume = dataToInt64(option.Data)
		case "itemCon":
			skill.ItemCon = dataToInt64(option.Data)
		case "itemConNo":
			skill.ItemConNo = dataToInt64(option.Data)
		case "mastery":
			skill.Mastery = dataToInt64(option.Data)
		case "time":
			skill.Time = dataToInt64(option.Data)
		case "eva":
			skill.Eva = dataToInt64(option.Data)
		case "acc":
			skill.Acc = dataToInt64(option.Data)
		case "jump":
			skill.Jump = dataToInt64(option.Data)
		case "speed":
			skill.Speed = dataToInt64(option.Data)
		case "range":
			skill.Range = dataToInt64(option.Data)
		case "mobCount":
			skill.MobCount = dataToInt64(option.Data)
		case "attackCount":
			skill.AttackCount = dataToInt64(option.Data)
		case "damage":
			skill.Damage = dataToInt64(option.Data)
		case "fixdamage":
			skill.Fixdamage = dataToInt64(option.Data)
		case "rb":
			skill.Rb = dataToVector(option.Data)
		case "hs":
			skill.Hs = textLookup[dataToUint32(option.Data)]
		case "lt":
			skill.Lt = dataToVector(option.Data)
		case "x":
			skill.X = dataToInt64(option.Data)
		case "y":
			skill.Y = dataToInt64(option.Data)
		case "z":
			skill.Z = dataToInt64(option.Data)
		case "prop":
			skill.Prop = dataToInt64(option.Data)
		case "ball":
		case "hit":
		case "bulletCount":
			skill.BulletCount = dataToInt64(option.Data)
		case "action":
			skill.Action = textLookup[dataToUint32(option.Data)]
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
