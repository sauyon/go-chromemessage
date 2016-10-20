package chromemsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"unsafe"
)

var nativeEndian binary.ByteOrder = endianness()

var defaultMsgr = Messenger{bufio.NewReadWriter(
	bufio.NewReader(os.Stdin),
	bufio.NewWriter(os.Stdout))}

type Messenger struct {
	port *bufio.ReadWriter
}

func New(port *bufio.ReadWriter) *Messenger {
	return &Messenger{port}
}

func Read(data interface{}) {
	defaultMsgr.Read(data)
}

func Write(msg interface{}) error {
	return defaultMsgr.Write(msg)
}

func (msgr *Messenger) Read(data interface{}) {
	lengthBits := make([]byte, 4)
	msgr.port.Read(lengthBits)
	length := nativeToInt(lengthBits)
	content := make([]byte, length)
	msgr.port.Read(content)
	json.Unmarshal(content, data)
}

func (msgr *Messenger) Write(msg interface{}) error {
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	length := len(json)
	bits := make([]byte, 4)
	buf := bytes.NewBuffer(bits)
	err = binary.Write(buf, nativeEndian, length)
	if err != nil {
		return err
	}
	_, err = msgr.port.Write(bits)
	if err != nil {
		return err
	}

	_, err = msgr.port.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func nativeToInt(bits []byte) int {
	var length uint32
	buf := bytes.NewBuffer(bits)
	binary.Read(buf, nativeEndian, &length)
	return int(length)
}

func endianness() binary.ByteOrder {
	var i int = 1
	bs := (*[unsafe.Sizeof(0)]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}
