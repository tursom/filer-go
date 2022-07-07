package filer

import (
	"syscall/js"
)

var (
	Console = ConsoleUtil(js.Global().Get("console"))
	Filer   = js.Global().Get("window").Get("Filer")
)

type (
	ConsoleUtil js.Value

	File struct {
		path string
	}
)

func (j ConsoleUtil) JsValue() js.Value {
	return js.Value(j)
}

func (c ConsoleUtil) Log(values ...any) {
	js.Value(c).Call("log", values...)
}

func Open(path string) *File {
	return &File{
		path: path,
	}
}

func (f *File) Write(p []byte) (n int, err error) {
	return len(p), Fs.AppendFile(f.path, p)
}

func (f *File) Read(p []byte) (n int, err error) {
	// TODO
	Fs.ReadFile(f.path)
	return
}
