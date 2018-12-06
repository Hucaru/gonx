package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Item data from nx
type Item struct {
	Cash            bool
	Unique          bool
	TradeBlock      bool
	ExpireOnLogout  bool
	Quest           bool
	InvTabID        byte
	MaxUpgradeSlots byte

	// Requirements
	Level                          byte
	JobID                          int16
	ReqStr, ReqDex, ReqInt, ReqLuk int16

	// Stat changes
	IncStr, IncDex, IncInt, IncLuk      int16
	IncAccuracy, IncEvasion             int16
	IncMagicDefence, IncPhysicalDefence int16
	IncMagicAttack, IncPhysicalAttack   int16
	IncMaxHP, IncMaxMP                  int16
	IncAttackSpeed, IncAttack           int16
	IncJump, IncSpeed                   int16

	// Shop information
	ShopPrice int32
	CanSell   bool

	// Pet
	Life, Hungry int16
}

// ExtractItems from parsed nx
func ExtractItems(nodes []Node, textLookup []string) map[int32]Item {
	items := make(map[int32]Item)

	searches := []string{"/Character/Accessory", "/Character/Cap", "/Character/Cape", "/Character/Coat",
		"/Character/Face", "/Character/Glove", "/Character/Hair", "/Character/Longcoat", "/Character/Pants",
		"/Character/PetEquip", "/Character/Ring", "/Character/Shield", "/Character/Shoes", "/Character/Weapon",
		"Item/Pet"}

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
					log.Println("Invalid node search:", subSearch)
				}

				name = strings.TrimSuffix(name, filepath.Ext(name))
				itemID, err := strconv.Atoi(name)

				if err != nil {
					log.Println(err)
					continue
				}

				item.InvTabID = byte(itemID / 1e6)
				items[int32(itemID)] = item
			}
		})

		if !valid {
			log.Println("Invalid node search:", search)
		}
	}

	searches = []string{"/Item/Cash", "/Item/Consume", "/Item/Etc", "/Item/Install", "/Item/Special"}

	for _, search := range searches {
		_ = search
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
			fallthrough
		case "incLUk":
			item.IncLuk = dataToInt16(option.Data)
		case "incMMD": // ?
			fallthrough
		case "incMDD":
			item.IncMagicDefence = dataToInt16(option.Data)
		case "incPDD":
			item.IncPhysicalDefence = dataToInt16(option.Data)
		case "incMAD":
			item.IncMagicAttack = dataToInt16(option.Data)
		case "incPAD":
			item.IncPhysicalAttack = dataToInt16(option.Data)
		case "incEVA":
			item.IncEvasion = dataToInt16(option.Data)
		case "incACC":
			item.IncAccuracy = dataToInt16(option.Data)
		case "incMHP":
			item.IncMaxHP = dataToInt16(option.Data)
		case "incMMP":
			item.IncMaxMP = dataToInt16(option.Data)
		case "only": // bool for only 1 of this item?
			item.Unique = dataToBool(option.Data[0])
		case "attackSpeed":
			item.IncAttackSpeed = dataToInt16(option.Data)
		case "attack":
			item.IncAttack = dataToInt16(option.Data)
		case "incSpeed":
			item.IncSpeed = dataToInt16(option.Data)
		case "incJump":
			item.IncJump = dataToInt16(option.Data)
		case "notSale":
			item.CanSell = dataToBool(option.Data[0])
		case "tradeBlock":
			item.TradeBlock = dataToBool(option.Data[0])
		case "expireOnLogout":
			item.ExpireOnLogout = dataToBool(option.Data[0])
		case "slotMax": // What do with this?
			// item.MaxUpgradeSlots = optionData[0]
		case "quest":
			item.Quest = dataToBool(option.Data[0])
		// I don't know what this is for
		case "tuc":
		case "timeLimited":
		case "recovery":
		case "reqPOP":
		case "regPOP":
		case "nameTag":
		case "pachinko":
		case "vslot":
		case "islot":

		// Not used
		case "icon":
		case "iconRaw":
		case "sfx":
		case "walk":
		case "afterImage":
		case "stand":
		case "knockback":
		case "fs":
		case "chatBalloon":
		case "sample":

		default:
			log.Println("Unsupported NX item option:", optionName)
		}

	}

	return item
}
