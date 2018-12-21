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
	ID      byte
	Pn      string
	Tm      int32
	Tn      string
	Pt      byte
	IsSpawn bool
	X, Y    int16
	Script  string
}

// Life object in a map
type Life struct {
	ID       int32
	IsMob    bool
	Foothold int16
	FaceLeft bool
	X, Y     int16
	MobTime  int64
	Hide     bool
	Rx0, Rx1 int32
	Cy       int32
	Info     byte
}

// Reactor object in a map
type Reactor struct {
	ID          int32
	FaceLeft    bool
	X, Y        int16
	ReactorTime int64
	Name        string
}

// Map data from nx
type Map struct {
	Town         bool
	ForcedReturn int32
	ReturnMap    int32
	MobRate      float64

	Swim, PersonalShop, EntrustedShop, ScrollDisable bool

	MoveLimit byte
	DecHP     int16

	NPCs     []Life
	Mobs     []Life
	Portals  []Portal
	Reactors []Reactor

	FieldLimit, VRLimit              byte
	VRRight, VRTop, VRLeft, VRBottom int16

	Recovery                  float64
	Version                   byte
	Bgm, MapMark              string
	Cloud, HideMinimap        bool
	MapDesc, Effect           string
	Fs                        float64
	TimeLimit                 int32
	FieldType                 byte
	Everlast, Snow, Rain      bool
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
			m.Town = dataToBool(option.Data[0])
		case "mobRate":
			m.MobRate = dataToFloat64(option.Data)
		case "forcedReturn":
			m.ForcedReturn = dataToInt32(option.Data)
		case "personalShop":
			m.PersonalShop = dataToBool(option.Data[0])
		case "entrustedShop":
			m.EntrustedShop = dataToBool(option.Data[0])
		case "swim":
			m.Swim = dataToBool(option.Data[0])
		case "moveLimit":
			m.MoveLimit = option.Data[0]
		case "decHP":
			m.DecHP = dataToInt16(option.Data)
		case "scrollDisable":
			m.ScrollDisable = dataToBool(option.Data[0])
		case "fieldLimit": // Max number of mobs on map?
			m.FieldLimit = option.Data[0]
		// Are VR settings to do with mob spawning? Determine which mob to spawn?
		case "VRRight":
			m.VRRight = dataToInt16(option.Data)
		case "VRTop":
			m.VRTop = dataToInt16(option.Data)
		case "VRLeft":
			m.VRLeft = dataToInt16(option.Data)
		case "VRBottom":
			m.VRBottom = dataToInt16(option.Data)
		case "VRLimit":
			m.VRLimit = option.Data[0]
		case "recovery": // float64
			m.Recovery = dataToFloat64(option.Data)
		case "returnMap":
			m.ReturnMap = dataToInt32(option.Data)
		case "version":
			m.Version = option.Data[0]
		case "bgm":
			m.Bgm = textLookup[dataToInt32(option.Data)]
		case "mapMark":
			m.MapMark = textLookup[dataToInt32(option.Data)]
		case "cloud":
			m.Cloud = dataToBool(option.Data[0])
		case "hideMinimap":
			m.HideMinimap = dataToBool(option.Data[0])
		case "mapDesc":
			m.MapDesc = textLookup[dataToInt32(option.Data)]
		case "effect":
			m.Effect = textLookup[dataToInt32(option.Data)]
		case "fs":
			m.Fs = dataToFloat64(option.Data)
		case "timeLimit": // is this for maps where a user can only be in there for x time?
			m.TimeLimit = dataToInt32(option.Data)
		case "fieldType":
			m.FieldType = option.Data[0]
		case "everlast":
			m.Everlast = dataToBool(option.Data[0])
		case "snow":
			m.Snow = dataToBool(option.Data[0])
		case "rain":
			m.Rain = dataToBool(option.Data[0])
		case "mapName":
			m.MapName = textLookup[dataToInt32(option.Data)]
		case "streetName":
			m.StreetName = textLookup[dataToInt32(option.Data)]
		case "help":
			m.Help = textLookup[dataToInt32(option.Data)]
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

		portal := Portal{ID: byte(portalNumber)}

		for j := uint32(0); j < uint32(portalObj.ChildCount); j++ {
			option := nodes[portalObj.ChildID+j]
			optionName := textLookup[option.NameID]

			switch optionName {
			case "pt":
				portal.Pt = option.Data[0]
			case "pn":
				portal.IsSpawn = bool(textLookup[dataToInt32(option.Data)] == "sp")
				portal.Pn = textLookup[dataToInt32(option.Data)]
			case "tm":
				portal.Tm = dataToInt32(option.Data)
			case "tn":
				portal.Tn = textLookup[dataToInt32(option.Data)]
			case "x":
				portal.X = dataToInt16(option.Data)
			case "y":
				portal.Y = dataToInt16(option.Data)
			case "script":
				portal.Script = textLookup[dataToInt32(option.Data)]
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
				life.ID = dataToInt32(option.Data)
			case "type":
				life.IsMob = bool(textLookup[dataToUint32(option.Data)] == "m")
			case "fh":
				life.Foothold = dataToInt16(option.Data)
			case "f":
				life.FaceLeft = dataToBool(option.Data[0])
			case "x":
				life.X = dataToInt16(option.Data)
			case "y":
				life.Y = dataToInt16(option.Data)
			case "mobTime":
				life.MobTime = dataToInt64(option.Data)
			case "hide":
				life.Hide = dataToBool(option.Data[0])
			case "rx0":
				life.Rx0 = dataToInt32(option.Data)
			case "rx1":
				life.Rx1 = dataToInt32(option.Data)
			case "cy":
				life.Cy = dataToInt32(option.Data)
			case "info": // An npc in map 103000002.img has info field
				life.Info = option.Data[0]
			default:
				fmt.Println("Unsupported NX life option:", optionName, "->", option.Data)
			}
		}

		if life.IsMob {
			mobs = append(mobs, life)
		} else {
			npcs = append(npcs, life)
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
				reactor.ID = dataToInt32(option.Data)
			case "x":
				reactor.X = dataToInt16(option.Data)
			case "y":
				reactor.Y = dataToInt16(option.Data)
			case "f":
				reactor.FaceLeft = dataToBool(option.Data[0])
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
