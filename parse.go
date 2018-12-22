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

// GetBitmap data along with x,y information
func (n Node) GetBitmap(bitmaps [][]byte) ([]byte, uint16, uint16, error) {
	if n.Type != 5 {
		return nil, 0, 0, fmt.Errorf("Not a bitmap node")
	}

	id := DataToUint32(n.Data)

	if int(id) >= len(bitmaps) {
		return nil, 0, 0, fmt.Errorf("Bitmap ID index out of range")
	}

	x := [8]byte{n.Data[4], n.Data[5]}
	y := [8]byte{n.Data[6], n.Data[7]}

	return bitmaps[id], DataToUint16(x), DataToUint16(y), nil
}

// GetAudio from the audio corresponding to this node
func (n Node) GetAudio(audio [][]byte) ([]byte, error) {
	if n.Type != 6 {
		return nil, fmt.Errorf("Not an audio node")
	}

	id := DataToUint32(n.Data)

	if int(id) >= len(audio) {
		return nil, fmt.Errorf("Audio ID index out of range")
	}

	return audio[id], nil
}

const packedNodeSize = 20 // size of packed Node struct

// Parse the nx file
func Parse(fname string) ([]Node, []string, [][]byte, [][]byte, error) {
	fBytes, err := ioutil.ReadFile(fname)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	buff := bytes.NewReader(fBytes)
	head := header{}
	err = binary.Read(buff, binary.LittleEndian, &head)

	if err != nil {
		return nil, nil, nil, nil, err
	}

	if head.Magic != [4]byte{0x50, 0x4B, 0x47, 0x34} {
		return nil, nil, nil, nil, fmt.Errorf("Not valid nx magic number")
	}

	nodes, err := readNodes(buff, head)
	strings, err := readStrings(buff, head)
	bitmaps, err := readBitmaps(buff, head)
	audio, err := readAudio(buff, head, nodes)

	return nodes, strings, bitmaps, audio, err
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

func readBitmaps(f *bytes.Reader, head header) ([][]byte, error) {
	_, err := f.Seek(head.BitmapOffsetTableOffset, 0)

	if err != nil {
		return nil, err
	}

	bitmapOffsets := make([]int64, head.BitmapCount)
	err = binary.Read(f, binary.LittleEndian, &bitmapOffsets)

	if err != nil {
		return nil, err
	}

	bitmapLookup := make([][]byte, head.BitmapCount)

	for i, v := range bitmapOffsets {
		_, err = f.Seek(v, 0)

		if err != nil {
			return nil, err
		}

		var length uint16
		err = binary.Read(f, binary.LittleEndian, &length)

		if err != nil {
			return nil, err
		}

		img := make([]byte, length)
		_, err = f.Read(img)

		if err != nil {
			return nil, err
		}

		bitmapLookup[i] = img
	}

	return bitmapLookup, nil
}

func readAudio(f *bytes.Reader, head header, nodes []Node) ([][]byte, error) {
	_, err := f.Seek(head.AudioOffsetTableOffset, 0)

	if err != nil {
		return nil, err
	}

	audioOffsets := make([]int64, head.AudioCount)
	err = binary.Read(f, binary.LittleEndian, &audioOffsets)

	if err != nil {
		return nil, err
	}

	audioLookup := make([][]byte, head.AudioCount)

	if head.AudioCount == 0 {
		return audioLookup, nil
	}

	for _, n := range nodes {
		if n.Type == 6 {
			id := DataToUint32(n.Data)

			offset := audioOffsets[id]

			_, err = f.Seek(offset, 0)

			if err != nil {
				return nil, err
			}

			tmp := [8]byte{n.Data[3], n.Data[4], n.Data[5], n.Data[6], n.Data[7]}
			length := DataToUint16(tmp)

			audio := make([]byte, length)
			_, err = f.Read(audio)

			if err != nil {
				return nil, err
			}

			audioLookup[id] = audio
		}
	}

	return audioLookup, nil
}

func readNodes(f *bytes.Reader, head header) ([]Node, error) {
	_, err := f.Seek(head.NodeBlockOffset, 0)

	if err != nil {
		return nil, err
	}

	if unsafe.Sizeof(Node{}) == packedNodeSize {
		nodeBytes := make([]byte, packedNodeSize*head.NodeCount)

		_, err = f.Read(nodeBytes)

		if err != nil {
			panic(err)
		}

		nodes := *(*[]Node)(unsafe.Pointer(&nodeBytes))
		nodes = nodes[:head.NodeCount] // Update slice header to have correct length

		if len(nodes) != int(head.NodeCount) {
			fmt.Println("Error in c style casting of nodes. Falling back to slower method")
			f.Seek(head.NodeBlockOffset, 0)
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

// FindNode with certain name e.g. element/sublement/grandchild
func FindNode(search string, nodes []Node, textLookup []string, fnc func(*Node)) bool {
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
