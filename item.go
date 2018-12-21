package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Item data from nx
type Item struct {
	Cash, Only, TradeBlock, ExpireOnLogout, Quest, TimeLimited     bool
	InvTabID, ReqLevel                                             byte
	Tuc                                                            byte // Total upgrade count?
	SlotMax                                                        int16
	ReqJob                                                         int16
	ReqSTR, ReqDEX, ReqINT, ReqLUK, IncSTR, IncDEX, IncINT, IncLUK int16
	IncACC, IncEVA, IncMDD, IncPDD, IncMAD, IncPAD, IncMHP, IncMMP int16
	AttackSpeed, Attack, IncJump, IncSpeed, RecoveryHP             int16
	Price                                                          int32
	NotSale                                                        bool
	UnitPrice                                                      float64
	Life, Hungry                                                   int16
	PickupItem, PickupAll, SweepForDrop                            bool
	ConsumeHP, LongRange                                           bool
	Recovery                                                       float64
	ReqPOP                                                         byte // ?
	NameTag                                                        byte
	Pachinko                                                       bool
	VSlot, ISlot                                                   string
	Type                                                           byte
	Success                                                        byte // Scroll type
	Cursed                                                         byte
	Add                                                            bool // ?
	DropSweep                                                      bool
	Rate                                                           byte
	Meso                                                           int32
	Path                                                           string
	FloatType                                                      byte
	NoFlip                                                         string
	StateChangeItem                                                int32
	BigSize                                                        byte
	Sfx                                                            string
	Walk                                                           byte
	AfterImage                                                     string
	Stand                                                          byte
	Knockback                                                      byte
	Fs                                                             byte
	ChatBalloon                                                    byte
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
			item.Cash = dataToBool(option.Data[0])
		case "reqSTR":
			item.ReqSTR = dataToInt16(option.Data)
		case "reqDEX":
			item.ReqDEX = dataToInt16(option.Data)
		case "reqINT":
			item.ReqINT = dataToInt16(option.Data)
		case "reqLUK":
			item.ReqLUK = dataToInt16(option.Data)
		case "reqJob":
			item.ReqJob = dataToInt16(option.Data)
		case "reqLevel":
			item.ReqLevel = option.Data[0]
		case "price":
			item.Price = dataToInt32(option.Data)
		case "incSTR":
			item.IncSTR = dataToInt16(option.Data)
		case "incDEX":
			item.IncDEX = dataToInt16(option.Data)
		case "incINT":
			item.IncINT = dataToInt16(option.Data)
		case "incLUK": // typo?
			fallthrough
		case "incLUk":
			item.IncLUK = dataToInt16(option.Data)
		case "incMMD": // typo?
			fallthrough
		case "incMDD":
			item.IncMDD = dataToInt16(option.Data)
		case "incPDD":
			item.IncPDD = dataToInt16(option.Data)
		case "incMAD":
			item.IncMAD = dataToInt16(option.Data)
		case "incPAD":
			item.IncPAD = dataToInt16(option.Data)
		case "incEVA":
			item.IncEVA = dataToInt16(option.Data)
		case "incACC":
			item.IncACC = dataToInt16(option.Data)
		case "incMHP":
			item.IncMHP = dataToInt16(option.Data)
		case "recoveryHP":
			item.RecoveryHP = dataToInt16(option.Data)
		case "incMMP":
			item.IncMMP = dataToInt16(option.Data)
		case "only":
			item.Only = dataToBool(option.Data[0])
		case "attackSpeed":
			item.AttackSpeed = dataToInt16(option.Data)
		case "attack":
			item.Attack = dataToInt16(option.Data)
		case "incSpeed":
			item.IncSpeed = dataToInt16(option.Data)
		case "incJump":
			item.IncJump = dataToInt16(option.Data)
		case "tuc": // total upgrade count?
			item.Tuc = option.Data[0]
		case "notSale":
			item.NotSale = dataToBool(option.Data[0])
		case "tradeBlock":
			item.TradeBlock = dataToBool(option.Data[0])
		case "expireOnLogout":
			item.ExpireOnLogout = dataToBool(option.Data[0])
		case "slotMax":
			item.SlotMax = dataToInt16(option.Data)
		case "quest":
			item.Quest = dataToBool(option.Data[0])
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
			item.UnitPrice = dataToFloat64(option.Data)
		case "timeLimited":
			item.TimeLimited = dataToBool(option.Data[0])
		case "recovery":
			item.Recovery = dataToFloat64(option.Data)
		case "regPOP":
			fallthrough
		case "reqPOP":
			item.ReqPOP = option.Data[0]
		case "nameTag":
			item.NameTag = option.Data[0]
		case "pachinko":
			item.Pachinko = dataToBool(option.Data[0])
		case "vslot":
			item.VSlot = textLookup[dataToInt32(option.Data)]
		case "islot":
			item.ISlot = textLookup[dataToInt32(option.Data)]
		case "type":
			item.Type = option.Data[0]
		case "success":
			item.Success = option.Data[0]
		case "cursed":
			item.Cursed = option.Data[0]
		case "add":
			item.Add = dataToBool(option.Data[0])
		case "dropSweep":
			item.DropSweep = dataToBool(option.Data[0])
		case "time":
		case "rate":
			item.Rate = option.Data[0]
		case "meso":
			item.Meso = dataToInt32(option.Data)
		case "path":
			idLookup := dataToUint32(option.Data)
			item.Path = textLookup[idLookup]
		case "floatType":
			item.FloatType = option.Data[0]
		case "noFlip":
			item.NoFlip = textLookup[dataToInt32(option.Data)]
		case "stateChangeItem":
			item.StateChangeItem = dataToInt32(option.Data)
		case "bigSize":
			item.BigSize = option.Data[0]
		case "icon":
		case "iconRaw":
		case "sfx":
			item.Sfx = textLookup[dataToInt32(option.Data)]
		case "walk":
			item.Walk = option.Data[0]
		case "afterImage":
			item.AfterImage = textLookup[dataToInt32(option.Data)]
		case "stand":
			item.Stand = option.Data[0]
		case "knockback":
			item.Knockback = option.Data[0]
		case "fs":
			item.Fs = option.Data[0]
		case "chatBalloon":
			item.ChatBalloon = option.Data[0]
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
