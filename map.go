package gonx

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

// Portal object in a map
type Portal struct {
	ID     int64
	Pn     string
	Tm     int64
	Tn     string
	Pt     int64
	X, Y   int64
	Script string
}

// Life object in a map
type Life struct {
	ID       int64
	Type     string
	Foothold int64
	FaceLeft int64
	X, Y     int64
	MobTime  int64
	Hide     int64
	Rx0, Rx1 int64
	Cy       int64
	Info     int64
}

// Reactor object in a map
type Reactor struct {
	ID          int64
	FaceLeft    int64
	X, Y        int64
	ReactorTime int64
	Name        string
}

// Map data from nx
type Map struct {
	Town         int64
	ForcedReturn int64
	ReturnMap    int64
	MobRate      float64

	Swim, PersonalShop, EntrustedShop, ScrollDisable int64

	MoveLimit int64
	DecHP     int64

	NPCs     []Life
	Mobs     []Life
	Portals  []Portal
	Reactors []Reactor

	FieldLimit, VRLimit              int64
	VRRight, VRTop, VRLeft, VRBottom int64

	Recovery                  float64
	Version                   int64
	Bgm, MapMark              string
	Cloud, HideMinimap        int64
	MapDesc, Effect           string
	Fs                        float64
	TimeLimit                 int64
	FieldType                 int64
	Everlast, Snow, Rain      int64
	MapName, StreetName, Help string
}

// ExtractMaps from parsed nx
func ExtractMaps(nodes []Node, textLookup []string) map[int32]Map {
	maps := make(map[int32]Map)

	searches := []string{"/Map/Map/Map0", "/Map/Map/Map1", "/Map/Map/Map2", "/Map/Map/Map9"}

	for _, search := range searches {
		valid := searchNode(search, nodes, textLookup, func(node *Node) {
			for i := uint32(0); i < uint32(node.ChildCount); i++ {
				mapNode := nodes[node.ChildID+i]
				name := textLookup[mapNode.NameID]

				var mapItem Map

				valid := searchNode(search+"/"+name+"/info", nodes, textLookup, func(node *Node) {
					mapItem = getMapInfo(node, nodes, textLookup)
				})

				if !valid {
					log.Println("Invalid node search:", search)
				}

				searchNode(search+"/"+name+"/life", nodes, textLookup, func(node *Node) {
					mapItem.NPCs, mapItem.Mobs = getMapLifes(node, nodes, textLookup)
				})

				searchNode(search+"/"+name+"/portal", nodes, textLookup, func(node *Node) {
					mapItem.Portals = getMapPortals(node, nodes, textLookup)
				})

				searchNode(search+"/"+name+"/reactor", nodes, textLookup, func(node *Node) {
					mapItem.Reactors = getMapReactors(node, nodes, textLookup)
				})

				name = strings.TrimSuffix(name, filepath.Ext(name))
				mapID, err := strconv.Atoi(name)

				if err != nil {
					log.Println(err)
					continue
				}

				maps[int32(mapID)] = mapItem
			}
		})

		if !valid {
			log.Println("Invalid node search:", search)
		}

	}

	return maps
}

func getMapInfo(node *Node, nodes []Node, textLookup []string) Map {
	var m Map
	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		option := nodes[node.ChildID+i]
		optionName := textLookup[option.NameID]

		switch optionName {
		case "town":
			m.Town = dataToInt64(option.Data)
		case "mobRate":
			m.MobRate = dataToFloat64(option.Data)
		case "forcedReturn":
			m.ForcedReturn = dataToInt64(option.Data)
		case "personalShop":
			m.PersonalShop = dataToInt64(option.Data)
		case "entrustedShop":
			m.EntrustedShop = dataToInt64(option.Data)
		case "swim":
			m.Swim = dataToInt64(option.Data)
		case "moveLimit":
			m.MoveLimit = dataToInt64(option.Data)
		case "decHP":
			m.DecHP = dataToInt64(option.Data)
		case "scrollDisable":
			m.ScrollDisable = dataToInt64(option.Data)
		case "fieldLimit": // Max number of mobs on map?
			m.FieldLimit = dataToInt64(option.Data)
		// Are VR settings to do with mob spawning? Determine which mob to spawn?
		case "VRRight":
			m.VRRight = dataToInt64(option.Data)
		case "VRTop":
			m.VRTop = dataToInt64(option.Data)
		case "VRLeft":
			m.VRLeft = dataToInt64(option.Data)
		case "VRBottom":
			m.VRBottom = dataToInt64(option.Data)
		case "VRLimit":
			m.VRLimit = dataToInt64(option.Data)
		case "recovery": // float64
			m.Recovery = dataToFloat64(option.Data)
		case "returnMap":
			m.ReturnMap = dataToInt64(option.Data)
		case "version":
			m.Version = dataToInt64(option.Data)
		case "bgm":
			m.Bgm = textLookup[dataToUint32(option.Data)]
		case "mapMark":
			m.MapMark = textLookup[dataToUint32(option.Data)]
		case "cloud":
			m.Cloud = dataToInt64(option.Data)
		case "hideMinimap":
			m.HideMinimap = dataToInt64(option.Data)
		case "mapDesc":
			m.MapDesc = textLookup[dataToUint32(option.Data)]
		case "effect":
			m.Effect = textLookup[dataToUint32(option.Data)]
		case "fs":
			m.Fs = dataToFloat64(option.Data)
		case "timeLimit": // is this for maps where a user can only be in there for x time?
			m.TimeLimit = dataToInt64(option.Data)
		case "fieldType":
			m.FieldType = dataToInt64(option.Data)
		case "everlast":
			m.Everlast = dataToInt64(option.Data)
		case "snow":
			m.Snow = dataToInt64(option.Data)
		case "rain":
			m.Rain = dataToInt64(option.Data)
		case "mapName":
			m.MapName = textLookup[dataToUint32(option.Data)]
		case "streetName":
			m.StreetName = textLookup[dataToUint32(option.Data)]
		case "help":
			m.Help = textLookup[dataToUint32(option.Data)]
		default:
			log.Println("Unsupported NX map option:", optionName, "->", option.Data)
		}
	}

	return m
}

func getMapPortals(node *Node, nodes []Node, textLookup []string) []Portal {
	portal := make([]Portal, node.ChildCount)

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		portalObj := nodes[node.ChildID+i]

		portalNumber, err := strconv.Atoi(textLookup[portalObj.NameID])

		if err != nil {
			fmt.Println("Skiping portal as ID is not a number")
			continue
		}

		portal := Portal{ID: int64(portalNumber)}

		for j := uint32(0); j < uint32(portalObj.ChildCount); j++ {
			option := nodes[portalObj.ChildID+j]
			optionName := textLookup[option.NameID]

			switch optionName {
			case "pt":
				portal.Pt = dataToInt64(option.Data)
			case "pn":
				portal.Pn = textLookup[dataToUint32(option.Data)]
			case "tm":
				portal.Tm = dataToInt64(option.Data)
			case "tn":
				portal.Tn = textLookup[dataToUint32(option.Data)]
			case "x":
				portal.X = dataToInt64(option.Data)
			case "y":
				portal.Y = dataToInt64(option.Data)
			case "script":
				portal.Script = textLookup[dataToUint32(option.Data)]
			default:
				fmt.Println("Unsupported NX portal option:", optionName, "->", option.Data)
			}
		}
	}

	return portal
}

func getMapLifes(node *Node, nodes []Node, textLookup []string) ([]Life, []Life) {
	npcs := []Life{}
	mobs := []Life{}

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		lifeObj := nodes[node.ChildID+i]

		var life Life

		for j := uint32(0); j < uint32(lifeObj.ChildCount); j++ {
			option := nodes[lifeObj.ChildID+j]
			optionName := textLookup[option.NameID]

			switch optionName {
			case "id":
				life.ID = dataToInt64(option.Data)
			case "type":
				life.Type = textLookup[dataToUint32(option.Data)]
			case "fh":
				life.Foothold = dataToInt64(option.Data)
			case "f":
				life.FaceLeft = dataToInt64(option.Data)
			case "x":
				life.X = dataToInt64(option.Data)
			case "y":
				life.Y = dataToInt64(option.Data)
			case "mobTime":
				life.MobTime = dataToInt64(option.Data)
			case "hide":
				life.Hide = dataToInt64(option.Data)
			case "rx0":
				life.Rx0 = dataToInt64(option.Data)
			case "rx1":
				life.Rx1 = dataToInt64(option.Data)
			case "cy":
				life.Cy = dataToInt64(option.Data)
			case "info": // An npc in map 103000002.img has info field
				life.Info = dataToInt64(option.Data)
			default:
				fmt.Println("Unsupported NX life option:", optionName, "->", option.Data)
			}
		}

		if life.Type == "m" {
			mobs = append(mobs, life)
		} else if life.Type == "n" {
			npcs = append(npcs, life)
		} else {
			fmt.Println("Unsupported life type:", life.Type)
		}
	}

	return npcs, mobs
}

func getMapReactors(node *Node, nodes []Node, textLookup []string) []Reactor {
	reactors := make([]Reactor, node.ChildCount)

	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		reactorObj := nodes[node.ChildID+i]

		var reactor Reactor

		for j := uint32(0); j < uint32(reactorObj.ChildCount); j++ {
			option := nodes[reactorObj.ChildID+j]
			optionName := textLookup[option.NameID]

			switch optionName {
			case "id":
				reactor.ID = dataToInt64(option.Data)
			case "x":
				reactor.X = dataToInt64(option.Data)
			case "y":
				reactor.Y = dataToInt64(option.Data)
			case "f":
				reactor.FaceLeft = dataToInt64(option.Data)
			case "reactorTime":
				reactor.ReactorTime = dataToInt64(option.Data)
			case "name":
				reactor.Name = textLookup[dataToUint32(option.Data)] // boss, ludigate
			default:
				fmt.Println("Unsupported NX reactor option:", optionName, "->", option.Data)
			}
		}
	}

	return reactors
}
