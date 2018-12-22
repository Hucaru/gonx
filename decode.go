package gonx

import (
	"encoding/binary"
	"math"
)

// Vector type
type Vector struct {
	X, Y int32
}

// DataToVector the nx payload
func DataToVector(data [8]byte) Vector {
	v := Vector{}
	v.X = DataToInt32(data)
	v.Y = int32(data[4]) | int32(data[5])<<8 | int32(data[6])<<16 | int32(data[7])<<24
	return v
}

// DataToInt64 the nx payload
func DataToInt64(data [8]byte) int64 {
	return int64(data[0]) |
		int64(data[1])<<8 |
		int64(data[2])<<16 |
		int64(data[3])<<24 |
		int64(data[4])<<32 |
		int64(data[5])<<40 |
		int64(data[6])<<48 |
		int64(data[7])<<56
}

// DataToUint64 the nx payload
func DataToUint64(data [8]byte) uint64 {
	return uint64(data[0]) |
		uint64(data[1])<<8 |
		uint64(data[2])<<16 |
		uint64(data[3])<<24 |
		uint64(data[4])<<32 |
		uint64(data[5])<<40 |
		uint64(data[6])<<48 |
		uint64(data[7])<<56
}

// DataToUint32 the nx payload
func DataToUint32(data [8]byte) uint32 {
	return uint32(data[0]) |
		uint32(data[1])<<8 |
		uint32(data[2])<<16 |
		uint32(data[3])<<24
}

// DataToInt32 the nx payload
func DataToInt32(data [8]byte) int32 {
	return int32(data[0]) |
		int32(data[1])<<8 |
		int32(data[2])<<16 |
		int32(data[3])<<24
}

// DataToInt16 the nx payload
func DataToInt16(data [8]byte) int16 {
	return int16(data[0]) |
		int16(data[1])<<8
}

// DataToUint16 the nx payload
func DataToUint16(data [8]byte) uint16 {
	return uint16(data[0]) |
		uint16(data[1])<<8
}

// DataToBool the nx payload
func DataToBool(data byte) bool {
	if data == 1 {
		return true
	}

	return false
}

// DataToFloat64 the nx payload
func DataToFloat64(data [8]byte) float64 {
	bits := binary.LittleEndian.Uint64(data[:])
	return math.Float64frombits(bits)
}
