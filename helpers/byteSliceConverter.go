package helpers

import (
	"encoding/binary"
)

func ConvertIntToByteArray(nr uint16) [2]byte {
	byteArray := make([]byte, 2)
	binary.BigEndian.PutUint16(byteArray, nr)
	var tmp [2]byte
	copy(tmp[:], byteArray)

	return tmp
}

func ConvertByteArrayToInt(byteArray []byte) uint16 {
	return binary.BigEndian.Uint16(byteArray)
}