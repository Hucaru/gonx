package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Mob data from nx
type Mob struct {
	MaxHP, HPRecovery  int32
	MaxMP, MPRecovery  int32
	Level              byte
	Exp                int32
	MADamage, MDDamage int16
	PADamage, PDDamage int16
	Speed, Eva, Acc    int16
	SummonType         int8
	Boss, Undead       bool
	ElemAttr           string
	Link               int32
	FlySpeed           int32
	NoRegen            byte
	Invincible         bool
	SelfDestruction    int32
	ExplosiveReward    bool
	Skills             map[int16]byte
	Revives            []int32
	Fs                 float64
	Pushed             int32
	BodyAttack         bool
	NoFlip             bool
	NotAttack          bool
	FirstAttack        bool
	RemoveQuest        bool
	RemoveAfter        string
	PublicReward       bool
	HPTagBGColor       byte
	HPTagColor         byte
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
			mobID, err := strconv.Atoi(name)

			if err != nil {
				log.Println(err)
				continue
			}

			mobs[int32(mobID)] = mob
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
		case "maxHP":
			mob.MaxHP = dataToInt32(option.Data)
		case "hpRecovery":
			mob.HPRecovery = dataToInt32(option.Data)
		case "maxMP":
			mob.MaxMP = dataToInt32(option.Data)
		case "mpRecovery":
			mob.MPRecovery = dataToInt32(option.Data)
		case "level":
			mob.Level = option.Data[0]
		case "exp":
			mob.Exp = dataToInt32(option.Data)
		case "MADamage":
			mob.MADamage = dataToInt16(option.Data)
		case "MDDamage":
			mob.MDDamage = dataToInt16(option.Data)
		case "PADamage":
			mob.PADamage = dataToInt16(option.Data)
		case "PDDamage":
			mob.PDDamage = dataToInt16(option.Data)
		case "speed":
			mob.Speed = dataToInt16(option.Data)
		case "eva":
			mob.Eva = dataToInt16(option.Data)
		case "acc":
			mob.Acc = dataToInt16(option.Data)
		case "summonType":
			mob.SummonType = int8(option.Data[0])
		case "boss":
			mob.Boss = dataToBool(option.Data[0])
		case "undead":
			mob.Undead = dataToBool(option.Data[0])
		case "elemAttr":
			mob.ElemAttr = textLookup[dataToInt32(option.Data)]
		case "link":
			mob.Link = dataToInt32(option.Data)
		case "flySpeed":
			mob.FlySpeed = dataToInt32(option.Data)
		case "noregen": // is this for both hp/mp?
			mob.NoRegen = option.Data[0]
		case "invincible":
			mob.Invincible = dataToBool(option.Data[0])
		case "selfDestruction":
			mob.SelfDestruction = dataToInt32(option.Data)
		case "explosiveReward": // A way that mob drops can drop?
			mob.ExplosiveReward = dataToBool(option.Data[0])
		case "skill":
			mob.Skills = getSkills(&option, nodes, textLookup)
		case "revive":
			mob.Revives = getRevives(&option, nodes)
		case "fs":
			mob.Fs = dataToFloat64(option.Data)
		case "pushed":
			mob.Pushed = dataToInt32(option.Data)
		case "bodyAttack":
			mob.BodyAttack = dataToBool(option.Data[0])
		case "noFlip":
			mob.NoFlip = dataToBool(option.Data[0])
		case "notAttack":
			mob.NotAttack = dataToBool(option.Data[0])
		case "firstAttack":
			mob.FirstAttack = dataToBool(option.Data[0])
		case "removeQuest":
			mob.RemoveQuest = dataToBool(option.Data[0])
		case "removeAfter":
			idLookup := dataToUint32(option.Data)
			mob.RemoveAfter = textLookup[idLookup]
		case "publicReward":
			mob.PublicReward = dataToBool(option.Data[0])
		case "hpTagBgcolor":
			mob.HPTagBGColor = option.Data[0]
		case "hpTagColor":
			mob.HPTagColor = option.Data[0]
		default:
			log.Println("Unsupported NX mob option:", optionName, "->", option.Data)
		}
	}

	return mob
}

func getSkills(node *Node, nodes []Node, textLookup []string) map[int16]byte {
	skills := make(map[int16]byte)

	// need to subnode the children of the children to node
	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		skillDir := nodes[node.ChildID+i]

		var id int16
		var level byte

		for j := uint32(0); j < uint32(skillDir.ChildCount); j++ {
			option := nodes[skillDir.ChildID+j]
			optionName := textLookup[option.NameID]

			switch optionName {
			case "level":
				level = option.Data[0]
			case "skill":
				id = dataToInt16(option.Data)
			case "action":
			case "effectAfter":
			default:
				log.Println("Unsupported NX mob skill option:", optionName, "->", option.Data)
			}
		}

		skills[id] = level
	}

	return skills
}

func getRevives(node *Node, nodes []Node) []int32 {
	revives := make([]int32, node.ChildCount)

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		mobID := nodes[node.ChildID+i]
		revives[i] = dataToInt32(mobID.Data)
	}

	return revives
}
