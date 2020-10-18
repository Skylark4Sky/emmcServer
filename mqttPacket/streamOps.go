package mqttPacket

import (
	"encoding/binary"
	"io"
)

var endianFormat = binary.LittleEndian

func readBytes(byteBuf io.Reader, data interface{}) bool {
	err := binary.Read(byteBuf, endianFormat, data)
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
	if readBytes(byteBuf, &data) != true {
		data = nil
	}
	return
}

func writeBytes(byteBuf io.Writer, data interface{}) bool {
	err := binary.Write(byteBuf, endianFormat, data)
	if err != nil {
		return false
	}
	return true
}

func setUint8(byteBuf io.Writer,value uint8) bool {
	return writeBytes(byteBuf,value)
}

func setUint16(byteBuf io.Writer,value uint16) bool {
	return writeBytes(byteBuf,value)
}

func setUint32(byteBuf io.Writer,value uint32) bool {
	return writeBytes(byteBuf,value)
}

func setBtyes(byteBuf io.Writer,value []byte) bool {
	return writeBytes(byteBuf,value)
}
