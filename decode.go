package gonx

import (
	"encoding/binary"
	"math"
)

// Vector type
type Vector struct {
	X, Y int32
}

func dataToVector(data [8]byte) Vector {
	v := Vector{}
	v.X = dataToInt32(data)
	v.Y = int32(data[4]) | int32(data[5])<<8 | int32(data[6])<<16 | int32(data[7])<<24
	return v
}

func dataToInt64(data [8]byte) int64 {
	return int64(data[0]) |
		int64(data[1])<<8 |
		int64(data[2])<<16 |
		int64(data[3])<<24 |
		int64(data[4])<<32 |
		int64(data[5])<<40 |
		int64(data[6])<<48 |
		int64(data[7])<<56
}

func dataToUint64(data [8]byte) uint64 {
	return uint64(data[0]) |
		uint64(data[1])<<8 |
		uint64(data[2])<<16 |
		uint64(data[3])<<24 |
		uint64(data[4])<<32 |
		uint64(data[5])<<40 |
		uint64(data[6])<<48 |
		uint64(data[7])<<56
}

func dataToUint32(data [8]byte) uint32 {
	return uint32(data[0]) |
		uint32(data[1])<<8 |
		uint32(data[2])<<16 |
		uint32(data[3])<<24
}

func dataToInt32(data [8]byte) int32 {
	return int32(data[0]) |
		int32(data[1])<<8 |
		int32(data[2])<<16 |
		int32(data[3])<<24
}

func dataToInt16(data [8]byte) int16 {
	return int16(data[0]) |
		int16(data[1])<<8
}

func dataToUint16(data [8]byte) uint16 {
	return uint16(data[0]) |
		uint16(data[1])<<8
}

func dataToBool(data byte) bool {
	if data == 1 {
		return true
	}

	return false
}

func dataToFloat64(data [8]byte) float64 {
	bits := binary.LittleEndian.Uint64(data[:])
	return math.Float64frombits(bits)
}
