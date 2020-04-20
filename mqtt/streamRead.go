package mqtt

import (
	"encoding/binary"
	"encoding/hex"
	"io"
)

func readBytes(byteBuf io.Reader, data interface{}) bool {
	err := binary.Read(byteBuf, binary.LittleEndian, data)
	if err != nil {
		return false
	}
	return true
}

func getUint8(byteBuf io.Reader) (data uint8) {
	if readBytes(byteBuf, &data) != true {
		data = 0
	}
	return
}

func getUint16(byteBuf io.Reader) (data uint16) {
	if readBytes(byteBuf, &data) != true {
		data = 0
	}
	return
}

func getUint32(byteBuf io.Reader) (data uint32) {
	if readBytes(byteBuf, &data) != true {
		data = 0
	}
	return
}

func getBtyes(byteBuf io.Reader, length uint32) (data []byte) {
	data = make([]byte, length)
	err := binary.Read(byteBuf, binary.LittleEndian, &data)
	if err != nil {
		data = nil
	}
	return
}

func GetBinaryData(binaryData []byte) string {
	return hex.EncodeToString(binaryData)
}