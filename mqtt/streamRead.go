package Mqtt

import (
	"encoding/binary"
	"io"
)

func readBytes(byteBuf io.Reader, data interface{}) bool {
	err := binary.Read(byteBuf, binary.LittleEndian, data)
	if err != nil {
		return false
	}
	return true
}

func GetUint8(byteBuf io.Reader) (data uint8) {
	if readBytes(byteBuf, &data) != true {
		data = 0
	}
	return
}

func GetUint16(byteBuf io.Reader) (data uint16) {
	if readBytes(byteBuf, &data) != true {
		data = 0
	}
	return
}

func GetUint32(byteBuf io.Reader) (data uint32) {
	if readBytes(byteBuf, &data) != true {
		data = 0
	}
	return
}

func GetBtyes(byteBuf io.Reader, length uint32) (data []byte) {
	data = make([]byte, length)
	err := binary.Read(byteBuf, binary.LittleEndian, &data)
	if err != nil {
		data = nil
	}
	return
}