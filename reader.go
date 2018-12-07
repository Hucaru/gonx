package gonx

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
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

	nodes := make([]Node, head.NodeCount)
	err = binary.Read(f, binary.LittleEndian, &nodes)

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
