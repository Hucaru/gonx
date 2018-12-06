package gonx

import (
	"encoding/binary"
	"os"
	"strings"
)

type header struct {
	Magic                   [4]byte
	NodeCount               uint32
	NodeBlockOffset         int64
	StringCount             uint32
	StringOffsetTableOffset int64
	BitmapCount             uint32
	BitmapOffsetTableOffset int64
	AudioCount              uint32
	AudioOffsetTableOffset  int64
}

// Node in nx file
type Node struct {
	NameID     uint32
	ChildID    uint32
	ChildCount uint16
	Type       uint16
	Data       [8]byte
}

// Parse the nx file
func Parse(fname string) ([]Node, []string) {
	f, err := os.Open(fname)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	head := header{}
	err = binary.Read(f, binary.LittleEndian, &head)

	if err != nil {
		panic(err)
	}

	if head.Magic != [4]byte{0x50, 0x4B, 0x47, 0x34} {
		panic("Not valid nx magic number")
	}

	strings := readStrings(f, head)
	nodes := readNodes(f, head)

	return nodes, strings
}

func searchNode(search string, nodes []Node, textLookup []string, fnc func(*Node)) bool {
	cursor := &nodes[0]

	path := strings.Split(search, "/")

	if strings.Compare(path[0], "/") == 0 {
		path = path[1:]
	}

	for j, p := range path {
		for i := uint32(0); i < uint32(cursor.ChildCount); i++ {

			if cursor.ChildID+i > uint32(len(nodes))-1 {
				return false
			}

			if strings.Compare(textLookup[nodes[cursor.ChildID+i].NameID], p) == 0 {
				cursor = &nodes[cursor.ChildID+i]

				if j == len(path)-1 {
					fnc(cursor)
					return true
				}

				break
			}
		}
	}

	return false
}

func readStrings(f *os.File, head header) []string {
	_, err := f.Seek(head.StringOffsetTableOffset, 0)

	if err != nil {
		panic(err)
	}

	stringOffsets := make([]int64, head.StringCount)
	err = binary.Read(f, binary.LittleEndian, &stringOffsets)

	if err != nil {
		panic(err)
	}

	strLookup := make([]string, head.StringCount)

	for i, v := range stringOffsets {
		_, err = f.Seek(v, 0)

		if err != nil {
			panic(err)
		}

		var length uint16
		err = binary.Read(f, binary.LittleEndian, &length)

		if err != nil {
			panic(err)
		}

		str := make([]byte, length)
		_, err = f.Read(str)

		if err != nil {
			panic(err)
		}

		strLookup[i] = string(str)
	}

	return strLookup
}

func readNodes(f *os.File, head header) []Node {
	_, err := f.Seek(head.NodeBlockOffset, 0)

	if err != nil {
		panic(err)
	}

	nodes := make([]Node, head.NodeCount)
	err = binary.Read(f, binary.LittleEndian, &nodes)

	if err != nil {
		panic(err)
	}

	return nodes
}
