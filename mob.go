package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Mob data from nx
type Mob struct {
	MaxHP, HPRecovery  int64
	MaxMP, MPRecovery  int64
	Level              int64
	Exp                int64
	MADamage, MDDamage int64
	PADamage, PDDamage int64
	Speed, Eva, Acc    int64
	SummonType         int64
	Boss, Undead       int64
	ElemAttr           string
	Link               int64
	FlySpeed           int64
	NoRegen            int64
	Invincible         int64
	SelfDestruction    int64
	ExplosiveReward    int64
	Skills             map[int64]int64
	Revives            []int64
	Fs                 float64
	Pushed             int64
	BodyAttack         int64
	NoFlip             int64
	NotAttack          int64
	FirstAttack        int64
	RemoveQuest        int64
	RemoveAfter        string
	PublicReward       int64
	HPTagBGColor       int64
	HPTagColor         int64
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
			mob.MaxHP = dataToInt64(option.Data)
		case "hpRecovery":
			mob.HPRecovery = dataToInt64(option.Data)
		case "maxMP":
			mob.MaxMP = dataToInt64(option.Data)
		case "mpRecovery":
			mob.MPRecovery = dataToInt64(option.Data)
		case "level":
			mob.Level = dataToInt64(option.Data)
		case "exp":
			mob.Exp = dataToInt64(option.Data)
		case "MADamage":
			mob.MADamage = dataToInt64(option.Data)
		case "MDDamage":
			mob.MDDamage = dataToInt64(option.Data)
		case "PADamage":
			mob.PADamage = dataToInt64(option.Data)
		case "PDDamage":
			mob.PDDamage = dataToInt64(option.Data)
		case "speed":
			mob.Speed = dataToInt64(option.Data)
		case "eva":
			mob.Eva = dataToInt64(option.Data)
		case "acc":
			mob.Acc = dataToInt64(option.Data)
		case "summonType":
			mob.SummonType = dataToInt64(option.Data)
		case "boss":
			mob.Boss = dataToInt64(option.Data)
		case "undead":
			mob.Undead = dataToInt64(option.Data)
		case "elemAttr":
			mob.ElemAttr = textLookup[dataToUint32(option.Data)]
		case "link":
			mob.Link = dataToInt64(option.Data)
		case "flySpeed":
			mob.FlySpeed = dataToInt64(option.Data)
		case "noregen": // is this for both hp/mp?
			mob.NoRegen = dataToInt64(option.Data)
		case "invincible":
			mob.Invincible = dataToInt64(option.Data)
		case "selfDestruction":
			mob.SelfDestruction = dataToInt64(option.Data)
		case "explosiveReward": // A way that mob drops can drop?
			mob.ExplosiveReward = dataToInt64(option.Data)
		case "skill":
			mob.Skills = getSkills(&option, nodes, textLookup)
		case "revive":
			mob.Revives = getRevives(&option, nodes)
		case "fs":
			mob.Fs = dataToFloat64(option.Data)
		case "pushed":
			mob.Pushed = dataToInt64(option.Data)
		case "bodyAttack":
			mob.BodyAttack = dataToInt64(option.Data)
		case "noFlip":
			mob.NoFlip = dataToInt64(option.Data)
		case "notAttack":
			mob.NotAttack = dataToInt64(option.Data)
		case "firstAttack":
			mob.FirstAttack = dataToInt64(option.Data)
		case "removeQuest":
			mob.RemoveQuest = dataToInt64(option.Data)
		case "removeAfter":
			idLookup := dataToUint32(option.Data)
			mob.RemoveAfter = textLookup[idLookup]
		case "publicReward":
			mob.PublicReward = dataToInt64(option.Data)
		case "hpTagBgcolor":
			mob.HPTagBGColor = dataToInt64(option.Data)
		case "hpTagColor":
			mob.HPTagColor = dataToInt64(option.Data)
		default:
			log.Println("Unsupported NX mob option:", optionName, "->", option.Data)
		}
	}

	return mob
}

func getSkills(node *Node, nodes []Node, textLookup []string) map[int64]int64 {
	skills := make(map[int64]int64)

	// need to subnode the children of the children to node
	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		skillDir := nodes[node.ChildID+i]

		var id, level int64

		for j := uint32(0); j < uint32(skillDir.ChildCount); j++ {
			option := nodes[skillDir.ChildID+j]
			optionName := textLookup[option.NameID]

			switch optionName {
			case "level":
				level = dataToInt64(option.Data)
			case "skill":
				id = dataToInt64(option.Data)
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

func getRevives(node *Node, nodes []Node) []int64 {
	revives := make([]int64, node.ChildCount)

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		mobID := nodes[node.ChildID+i]
		revives[i] = dataToInt64(mobID.Data)
	}

	return revives
}
