package chromemsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
)

const nativeEndian binary.ByteOrder = endianness()

const defaultReader = MessageReader{bufio.NewReader(os.Stdin)}

type MessageReader struct {
	in *bufio.Reader
}

func New(in *bufio.Reader) *MessageReader {
	return MessageReader{in}
}

func Read(data *interface{}) {
	defaultReader.Read(data)
}

func (reader *MessageReader) Read(data *interface{}) {
	lengthBits := make([]byte, 4)
	reader.in.Read(lengthBits)
	length := nativeToInt(lengthBits)
	content := make([]byte, length)
	s.Read(content)
	json.Unmarshal(content, data)
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
