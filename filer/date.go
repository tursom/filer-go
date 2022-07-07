package filer

import "syscall/js"

type (
	JsDate js.Value
)

func (j JsDate) JsValue() js.Value {
	return js.Value(j)
}
