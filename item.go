package gonx

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Item data from nx
type Item struct {
	Cash   bool
	Unique bool
	Trade  bool
	SlotID byte

	// Requirements
	Level                          byte
	JobID                          int16
	ReqStr, ReqDex, ReqInt, ReqLuk int16

	// Stat changes
	IncStr, IncDex, IncInt, IncLuk, IncAccuracy, IncEvasion    int16
	MagicDefence, PhysicalDefence, MagicAttack, PhysicalAttack int16
	MaxHP, MaxMP                                               int16

	ShopPrice int32
}

// ExtractItems from parsed nx
func ExtractItems(nodes []Node, textLookup []string) map[int32]Item {
	items := make(map[int32]Item)

	searches := []string{"/Character/Accessory", "/Character/Cap", "/Character/Cape", "/Character/Coat",
		"/Character/Face", "/Character/Glove", "/Character/Hair", "/Character/Longvoat", "/Character/Pants",
		"/Character/PetEquip", "/Character/Ring", "/Character/Shield", "/Character/Shoes", "/Character/Weapon"}

	for _, search := range searches {
		valid := searchNode(search, nodes, textLookup, func(node *Node) {
			for i := uint32(0); i < uint32(node.ChildCount); i++ {
				mobID := nodes[node.ChildID+i]
				name := textLookup[mobID.NameID]
				subSearch := search + "/" + name + "/info"

				var item Item

				valid := searchNode(subSearch, nodes, textLookup, func(node *Node) {
					item = getItem(node, nodes, textLookup)
				})

				if !valid {
					log.Println("Invalid node search for:", subSearch)
				}

				name = strings.TrimSuffix(name, filepath.Ext(name))
				itemID, err := strconv.Atoi(name)

				if err != nil {
					log.Println(err)
				}

				items[int32(itemID)] = item
			}
		})

		if !valid {
			log.Println("Invalid node search:", search)
		}
	}

	return items
}

func getItem(node *Node, nodes []Node, textLookup []string) Item {
	item := Item{}

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		option := nodes[node.ChildID+i]
		optionName := textLookup[option.NameID]

		switch optionName {
		case "cash":
			item.Cash = dataToBool(option.Data[0])
		case "reqSTR":
			item.ReqStr = dataToInt16(option.Data)
		case "reqDEX":
			item.ReqDex = dataToInt16(option.Data)
		case "reqINT":
			item.ReqInt = dataToInt16(option.Data)
		case "reqLUK":
			item.ReqLuk = dataToInt16(option.Data)
		case "reqJob":
			item.JobID = dataToInt16(option.Data)
		case "reqLevel":
			item.Level = option.Data[0]
		case "price":
			item.ShopPrice = dataToInt32(option.Data)
		case "incSTR":
			item.IncStr = dataToInt16(option.Data)
		case "incDEX":
			item.IncDex = dataToInt16(option.Data)
		case "incINT":
			item.IncInt = dataToInt16(option.Data)
		case "incLUK":
			item.IncLuk = dataToInt16(option.Data)
		case "incMDD":
			item.MagicDefence = dataToInt16(option.Data)
		case "incPDD":
			item.PhysicalDefence = dataToInt16(option.Data)
		case "incMAD":
			item.MagicAttack = dataToInt16(option.Data)
		case "incPAD":
			item.PhysicalAttack = dataToInt16(option.Data)
		case "incEVA":
			item.IncEvasion = dataToInt16(option.Data)
		case "incACC":
			item.IncAccuracy = dataToInt16(option.Data)
		case "incMHP":
			item.MaxHP = dataToInt16(option.Data)
		case "incMMP":
			item.MaxMP = dataToInt16(option.Data)
		case "only": // bool for only 1 of this item?
			item.Unique = dataToBool(option.Data[0])
		case "tuc": // ?
		case "vslot": // visual slot?
		case "islot": // inventory slot?
			item.SlotID = option.Data[0]
		case "timeLimited":
		case "icon":
		case "iconRaw":

		default:
			fmt.Println(optionName)
		}

	}

	return item
}
