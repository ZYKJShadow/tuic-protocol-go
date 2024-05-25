package utils

import "encoding/binary"

func ReadUint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func WriteUint16(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}

func ReadUint8(b []byte) uint8 {
	return b[0]
}

func WriteUint8(b []byte, v uint8) {
	b[0] = v
}
