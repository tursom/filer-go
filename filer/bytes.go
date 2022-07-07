package filer

import (
	"syscall/js"
)

var (
	bytesArray = js.Global().Get("Uint8Array")
	BufferType = Filer.Get("Buffer")
)

type (
	Buffer js.Value
)

func (b Buffer) JsValue() js.Value {
	return js.Value(b)
}

func Bytes(bytes []byte) js.Value {
	buf := bytesArray.New(len(bytes))
	for i, b := range bytes {
		buf.SetIndex(i, b)
	}
	return buf
}

func WrapBuffer(buffer []byte) Buffer {
	return Buffer(BufferType.Call("from", Bytes(buffer)))
}

func AllocBuffer(len int) Buffer {
	return Buffer(BufferType.Call("alloc", len))
}

func (b Buffer) Length() int {
	return js.Value(b).Length()
}

func (b Buffer) GetIndex(i int) byte {
	return byte(js.Value(b).Index(i).Int())
}

func (b Buffer) SetIndex(i int, x byte) {
	js.Value(b).SetIndex(i, x)
}

func (b Buffer) Write(bytes []byte, offset, length, position int) {
	for i := 0; i < length; i++ {
		b.SetIndex(i+position, bytes[offset+i])
	}
}

func (b Buffer) Read(bytes []byte, offset, length, position int) {
	for i := 0; i < length; i++ {
		bytes[offset+i] = b.GetIndex(i + position)
	}
}
