package gonx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"strings"
	"unsafe"
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

const packedNodeSize = 20 // size of packed Node struct

// Parse the nx file
func Parse(fname string) ([]Node, []string, error) {
	fBytes, err := ioutil.ReadFile(fname)

	if err != nil {
		return nil, nil, err
	}

	buff := bytes.NewReader(fBytes)
	head := header{}
	err = binary.Read(buff, binary.LittleEndian, &head)

	if err != nil {
		return nil, nil, err
	}

	if head.Magic != [4]byte{0x50, 0x4B, 0x47, 0x34} {
		return nil, nil, fmt.Errorf("Not valid nx magic number")
	}

	strings, err := readStrings(buff, head)
	nodes, err := readNodes(buff, head)

	return nodes, strings, err
}

func readStrings(f *bytes.Reader, head header) ([]string, error) {
	_, err := f.Seek(head.StringOffsetTableOffset, 0)

	if err != nil {
		return nil, err
	}

	stringOffsets := make([]int64, head.StringCount)
	err = binary.Read(f, binary.LittleEndian, &stringOffsets)

	if err != nil {
		return nil, err
	}

	strLookup := make([]string, head.StringCount)

	for i, v := range stringOffsets {
		_, err = f.Seek(v, 0)

		if err != nil {
			return nil, err
		}

		var length uint16
		err = binary.Read(f, binary.LittleEndian, &length)

		if err != nil {
			return nil, err
		}

		str := make([]byte, length)
		_, err = f.Read(str)

		if err != nil {
			return nil, err
		}

		strLookup[i] = string(str)
	}

	return strLookup, nil
}

func readNodes(f *bytes.Reader, head header) ([]Node, error) {
	_, err := f.Seek(head.NodeBlockOffset, 0)

	if err != nil {
		return nil, err
	}

	if unsafe.Sizeof(Node{}) == packedNodeSize {
		nodeBytes := make([]byte, packedNodeSize*head.NodeCount)
		f.Read(nodeBytes)
		nodes := *(*[]Node)(unsafe.Pointer(&nodeBytes))
		if len(nodes) != int(head.NodeCount*packedNodeSize) {
			fmt.Println("Error in c style casting of nodes. Falling back to slower method")
			goto FALLBACK
		}

		return nodes, nil
	FALLBACK:
	}
	nodes := make([]Node, head.NodeCount)

	for i := range nodes {
		nameID := make([]byte, 4)
		childID := make([]byte, 4)
		childCount := make([]byte, 2)
		ntype := make([]byte, 2)
		data := make([]byte, 8)

		_, err = f.Read(nameID)

		if err != nil {
			return nil, err
		}

		_, err = f.Read(childID)

		if err != nil {
			return nil, err
		}

		_, err = f.Read(childCount)

		if err != nil {
			return nil, err
		}

		_, err = f.Read(ntype)

		if err != nil {
			return nil, err
		}

		_, err = f.Read(data)

		if err != nil {
			return nil, err
		}

		nodes[i].NameID = binary.LittleEndian.Uint32(nameID)
		nodes[i].ChildID = binary.LittleEndian.Uint32(childID)
		nodes[i].ChildCount = binary.LittleEndian.Uint16(childCount)
		nodes[i].Type = binary.LittleEndian.Uint16(ntype)

		// Fixed size and small. Unroll the loop
		nodes[i].Data[7] = data[7] // slice length check optimisation
		nodes[i].Data[0] = data[0]
		nodes[i].Data[1] = data[1]
		nodes[i].Data[2] = data[2]
		nodes[i].Data[3] = data[3]
		nodes[i].Data[4] = data[4]
		nodes[i].Data[5] = data[5]
		nodes[i].Data[6] = data[6]
	}

	if err != nil {
		return nil, err
	}

	return nodes, nil
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
