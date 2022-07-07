package filer

import "syscall/js"

type JsErr js.Value

func (j JsErr) JsValue() js.Value {
	return js.Value(j)
}

func (j JsErr) Error() string {
	return js.Value(j).Get("message").String()
}
