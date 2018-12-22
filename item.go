package gonx

import (
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Item data from nx
type Item struct {
	InvTabID                                                       byte
	Cash, Only, TradeBlock, ExpireOnLogout, Quest, TimeLimited     int64
	ReqLevel                                                       int64
	Tuc                                                            int64 // Total upgrade count?
	SlotMax                                                        int64
	ReqJob                                                         int64
	ReqSTR, ReqDEX, ReqINT, ReqLUK, IncSTR, IncDEX, IncINT, IncLUK int64
	IncACC, IncEVA, IncMDD, IncPDD, IncMAD, IncPAD, IncMHP, IncMMP int64
	AttackSpeed, Attack, IncJump, IncSpeed, RecoveryHP             int64
	Price                                                          int64
	NotSale                                                        int64
	UnitPrice                                                      float64
	Life, Hungry                                                   int64
	PickupItem, PickupAll, SweepForDrop                            int64
	ConsumeHP, LongRange                                           int64
	Recovery                                                       float64
	ReqPOP                                                         int64 // ?
	NameTag                                                        int64
	Pachinko                                                       int64
	VSlot, ISlot                                                   string
	Type                                                           int64
	Success                                                        int64 // Scroll type
	Cursed                                                         int64
	Add                                                            int64 // ?
	DropSweep                                                      int64
	Rate                                                           int64
	Meso                                                           int64
	Path                                                           string
	FloatType                                                      int64
	NoFlip                                                         string
	StateChangeItem                                                int64
	BigSize                                                        int64
	Sfx                                                            string
	Walk                                                           int64
	AfterImage                                                     string
	Stand                                                          int64
	Knockback                                                      int64
	Fs                                                             int64
	ChatBalloon                                                    int64
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
			item.Cash = dataToInt64(option.Data)
		case "reqSTR":
			item.ReqSTR = dataToInt64(option.Data)
		case "reqDEX":
			item.ReqDEX = dataToInt64(option.Data)
		case "reqINT":
			item.ReqINT = dataToInt64(option.Data)
		case "reqLUK":
			item.ReqLUK = dataToInt64(option.Data)
		case "reqJob":
			item.ReqJob = dataToInt64(option.Data)
		case "reqLevel":
			item.ReqLevel = dataToInt64(option.Data)
		case "price":
			item.Price = dataToInt64(option.Data)
		case "incSTR":
			item.IncSTR = dataToInt64(option.Data)
		case "incDEX":
			item.IncDEX = dataToInt64(option.Data)
		case "incINT":
			item.IncINT = dataToInt64(option.Data)
		case "incLUK": // typo?
			fallthrough
		case "incLUk":
			item.IncLUK = dataToInt64(option.Data)
		case "incMMD": // typo?
			fallthrough
		case "incMDD":
			item.IncMDD = dataToInt64(option.Data)
		case "incPDD":
			item.IncPDD = dataToInt64(option.Data)
		case "incMAD":
			item.IncMAD = dataToInt64(option.Data)
		case "incPAD":
			item.IncPAD = dataToInt64(option.Data)
		case "incEVA":
			item.IncEVA = dataToInt64(option.Data)
		case "incACC":
			item.IncACC = dataToInt64(option.Data)
		case "incMHP":
			item.IncMHP = dataToInt64(option.Data)
		case "recoveryHP":
			item.RecoveryHP = dataToInt64(option.Data)
		case "incMMP":
			item.IncMMP = dataToInt64(option.Data)
		case "only":
			item.Only = dataToInt64(option.Data)
		case "attackSpeed":
			item.AttackSpeed = dataToInt64(option.Data)
		case "attack":
			item.Attack = dataToInt64(option.Data)
		case "incSpeed":
			item.IncSpeed = dataToInt64(option.Data)
		case "incJump":
			item.IncJump = dataToInt64(option.Data)
		case "tuc": // total upgrade count?
			item.Tuc = dataToInt64(option.Data)
		case "notSale":
			item.NotSale = dataToInt64(option.Data)
		case "tradeBlock":
			item.TradeBlock = dataToInt64(option.Data)
		case "expireOnLogout":
			item.ExpireOnLogout = dataToInt64(option.Data)
		case "slotMax":
			item.SlotMax = dataToInt64(option.Data)
		case "quest":
			item.Quest = dataToInt64(option.Data)
		case "life":
			item.Life = dataToInt64(option.Data)
		case "hungry":
			item.Hungry = dataToInt64(option.Data)
		case "pickupItem":
			item.PickupItem = dataToInt64(option.Data)
		case "pickupAll":
			item.PickupAll = dataToInt64(option.Data)
		case "sweepForDrop":
			item.SweepForDrop = dataToInt64(option.Data)
		case "longRange":
			item.LongRange = dataToInt64(option.Data)
		case "consumeHP":
			item.ConsumeHP = dataToInt64(option.Data)
		case "unitPrice":
			item.UnitPrice = dataToFloat64(option.Data)
		case "timeLimited":
			item.TimeLimited = dataToInt64(option.Data)
		case "recovery":
			item.Recovery = dataToFloat64(option.Data)
		case "regPOP":
			fallthrough
		case "reqPOP":
			item.ReqPOP = dataToInt64(option.Data)
		case "nameTag":
			item.NameTag = dataToInt64(option.Data)
		case "pachinko":
			item.Pachinko = dataToInt64(option.Data)
		case "vslot":
			item.VSlot = textLookup[dataToUint32(option.Data)]
		case "islot":
			item.ISlot = textLookup[dataToUint32(option.Data)]
		case "type":
			item.Type = dataToInt64(option.Data)
		case "success":
			item.Success = dataToInt64(option.Data)
		case "cursed":
			item.Cursed = dataToInt64(option.Data)
		case "add":
			item.Add = dataToInt64(option.Data)
		case "dropSweep":
			item.DropSweep = dataToInt64(option.Data)
		case "time":
		case "rate":
			item.Rate = dataToInt64(option.Data)
		case "meso":
			item.Meso = dataToInt64(option.Data)
		case "path":
			idLookup := dataToUint32(option.Data)
			item.Path = textLookup[idLookup]
		case "floatType":
			item.FloatType = dataToInt64(option.Data)
		case "noFlip":
			item.NoFlip = textLookup[dataToUint32(option.Data)]
		case "stateChangeItem":
			item.StateChangeItem = dataToInt64(option.Data)
		case "bigSize":
			item.BigSize = dataToInt64(option.Data)
		case "icon":
		case "iconRaw":
		case "sfx":
			item.Sfx = textLookup[dataToUint32(option.Data)]
		case "walk":
			item.Walk = dataToInt64(option.Data)
		case "afterImage":
			item.AfterImage = textLookup[dataToUint32(option.Data)]
		case "stand":
			item.Stand = dataToInt64(option.Data)
		case "knockback":
			item.Knockback = dataToInt64(option.Data)
		case "fs":
			item.Fs = dataToInt64(option.Data)
		case "chatBalloon":
			item.ChatBalloon = dataToInt64(option.Data)
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
