package gonx

import (
	"encoding/binary"
	"log"
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

// Item data from nx
type Item struct {
	// General
	IsCash          bool
	IsUnique        bool
	IsTradeBlock    bool
	ExpireOnLogout  bool
	IsQuest         bool
	InvTabID        byte
	MaxUpgradeCount byte
	MaxStackSize    int16

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
	RecoverHP                           int16

	// Shop information
	SellToShopPrice int32
	CanSell         bool
	UnitPrice       float64

	// Pet
	Life, Hungry                        int16
	PickupItem, PickupAll, SweepForDrop bool
	ConsumeHP, LongRange                bool
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
				itemNode := nodes[node.ChildID+i]
				name := textLookup[itemNode.NameID]
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

	searches = []string{"/Item/Cash", "/Item/Consume", "/Item/Etc", "/Item/Install"}

	for _, search := range searches {
		valid := searchNode(search, nodes, textLookup, func(node *Node) {
			for i := uint32(0); i < uint32(node.ChildCount); i++ {
				itemGroupNode := nodes[node.ChildID+i]
				groupName := textLookup[itemGroupNode.NameID]

				for j := uint32(0); j < uint32(itemGroupNode.ChildCount); j++ {
					itemNode := nodes[itemGroupNode.ChildID+j]
					name := textLookup[itemNode.NameID]

					subSearch := search + "/" + groupName + "/" + name + "/info"

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
			item.IsCash = dataToBool(option.Data[0])
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
			item.SellToShopPrice = dataToInt32(option.Data)
		case "incSTR":
			item.IncStr = dataToInt16(option.Data)
		case "incDEX":
			item.IncDex = dataToInt16(option.Data)
		case "incINT":
			item.IncInt = dataToInt16(option.Data)
		case "incLUK": // typo?
			fallthrough
		case "incLUk":
			item.IncLuk = dataToInt16(option.Data)
		case "incMMD": // typo?
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
		case "only":
			item.IsUnique = dataToBool(option.Data[0])
		case "attackSpeed":
			item.IncAttackSpeed = dataToInt16(option.Data)
		case "attack":
			item.IncAttack = dataToInt16(option.Data)
		case "incSpeed":
			item.IncSpeed = dataToInt16(option.Data)
		case "incJump":
			item.IncJump = dataToInt16(option.Data)
		case "tuc": // total upgrade count?
			item.MaxUpgradeCount = option.Data[0]
		case "notSale":
			item.CanSell = dataToBool(option.Data[0])
		case "tradeBlock":
			item.IsTradeBlock = dataToBool(option.Data[0])
		case "expireOnLogout":
			item.ExpireOnLogout = dataToBool(option.Data[0])
		case "slotMax":
			item.MaxStackSize = dataToInt16(option.Data)
		case "quest":
			item.IsQuest = dataToBool(option.Data[0])
		case "life":
			item.Life = dataToInt16(option.Data)
		case "hungry":
			item.Hungry = dataToInt16(option.Data)
		case "pickupItem":
			item.PickupItem = dataToBool(option.Data[0])
		case "pickupAll":
			item.PickupAll = dataToBool(option.Data[0])
		case "sweepForDrop":
			item.SweepForDrop = dataToBool(option.Data[0])
		case "longRange":
			item.LongRange = dataToBool(option.Data[0])
		case "consumeHP":
			item.ConsumeHP = dataToBool(option.Data[0])
		case "unitPrice":
			bits := binary.LittleEndian.Uint64(option.Data[:])
			item.UnitPrice = math.Float64frombits(bits)
		case "recoveryHP":
			item.RecoverHP = dataToInt16(option.Data)

		// I don't know what the following denote
		case "timeLimited":
		case "recovery": // float64
		case "reqPOP":
		case "regPOP":
		case "nameTag":
		case "pachinko":
		case "vslot":
		case "islot":
		case "type":
		case "success":
		case "cursed":
		case "add":
		case "dropSweep":
		case "time":
		case "rate":
		case "meso":
		case "path":
		case "floatType":
		case "noFlip":
		case "stateChangeItem":
		case "bigSize":

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
		case "iconD":
		case "iconRawD":
		case "iconReward":

		default:
			log.Println("Unsupported NX item option:", optionName, "->", option.Data)
		}

	}

	return item
}
