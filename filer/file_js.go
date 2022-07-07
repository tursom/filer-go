package filer

import (
	"encoding/hex"
	"sync"
	"syscall/js"
)

var (
	Fs = JsFs(Filer.Get("fs"))
)

type (
	JsFs        js.Value
	JsFd        js.Value
	AccessMode  js.Value
	WatchOption struct {
		// TODO
	}
	WatchEvent uint8
)

const (
	WatchEventUnsupported WatchEvent = 0
	WatchEventRename                 = 1
	WatchEventChange                 = 2
)

func (f JsFs) JsValue() js.Value {
	return js.Value(f)
}

func (f JsFd) JsValue() js.Value {
	return js.Value(f)
}

func (o WatchOption) JsValue() js.Value {
	return js.ValueOf(make(map[string]any))
}

func (m AccessMode) JsValue() js.Value {
	return js.Value(m)
}

func (f JsFs) call(
	m string,
	args ...any,
) (this js.Value, ret []js.Value) {
	var wait sync.WaitGroup
	wait.Add(1)

	args = append(args, js.FuncOf(func(callbackThis js.Value, r []js.Value) any {
		defer wait.Done()

		this = callbackThis
		ret = r

		return nil
	}))

	js.Value(Fs).Call(m, args...)

	wait.Wait()
	return
}

func hasError(ret []js.Value) bool {
	if len(ret) == 0 {
		return false
	}

	return ret[0].Type() != js.TypeNull
}

func wrapError(ret []js.Value) error {
	if hasError(ret) {
		return JsErr(ret[0])
	}

	return nil
}

func (f JsFs) Rename(oldPath, newPath string) (err error) {
	_, ret := f.call("rename", oldPath, newPath)
	return wrapError(ret)
}

func (f JsFs) Ftruncate(fd JsFd, len int) (err error) {
	_, ret := f.call("ftruncate", fd, len)
	return wrapError(ret)
}

func (f JsFs) Truncate(path string, len int) (err error) {
	_, ret := f.call("truncate", path, len)
	return wrapError(ret)
}

func (f JsFs) Stat(path string) (stats *FileStats, err error) {
	_, ret := f.call("stat", path)
	if hasError(ret) {
		return nil, JsErr(ret[0])
	}

	return GetFileStats(ret[1]), nil
}

func (f JsFs) Fstat(fd JsFd) (stats *FileStats, err error) {
	_, ret := f.call("fstat", fd)
	if hasError(ret) {
		return nil, JsErr(ret[0])
	}

	return GetFileStats(ret[1]), nil
}

func (f JsFs) Lstat(path string) (stats *FileStats, err error) {
	_, ret := f.call("lstat", path)
	if hasError(ret) {
		return nil, JsErr(ret[0])
	}

	return GetFileStats(ret[1]), nil
}

func (f JsFs) Exists(path string) bool {
	_, ret := f.call("exists", path)
	return ret[0].Bool()
}

func (f JsFs) Link(srcpath, dstpath string) (err error) {
	_, ret := f.call("link", srcpath, dstpath)
	return wrapError(ret)
}

func (f JsFs) Symlink(srcpath, dstpath string) (err error) {
	_, ret := f.call("symlink", srcpath, dstpath)
	return wrapError(ret)
}

func (f JsFs) Readlink(path string) (linkContents string, err error) {
	_, ret := f.call("readlink", path)
	if hasError(ret) {
		return "", JsErr(ret[0])
	}

	return ret[1].String(), nil
}

// Realpath Not implemented, see https://github.com/filerjs/filer/issues/85
func (f JsFs) Realpath(path string) (realpath string, err error) {
	_, ret := f.call("realpath", path)
	if hasError(ret) {
		return "", JsErr(ret[0])
	}

	return ret[1].String(), nil
}

func (f JsFs) Unlink(path string) (err error) {
	_, ret := f.call("unlink", path)
	return wrapError(ret)
}

func (f JsFs) Mknod(path, mode string) (err error) {
	_, ret := f.call("mknod", path, mode)
	return wrapError(ret)
}

func (f JsFs) Rmdir(path string) (err error) {
	_, ret := f.call("rmdir", path)
	return wrapError(ret)
}

func (f JsFs) Mkdir(path string) (err error) {
	_, ret := f.call("mkdir", path)
	return wrapError(ret)
}

func (f JsFs) Access(path, mode *AccessMode) (err error) {
	args := []any{path}
	if mode != nil {
		args = append(args, mode)
	}

	_, ret := f.call("acccess", args...)
	return wrapError(ret)
}

func (f JsFs) Mkdtemp(path string) (tmpPath string, err error) {
	_, ret := f.call("mkdtemp", path)
	if hasError(ret) {
		return "", JsErr(ret[0])
	}

	return ret[1].String(), nil
}

func (f JsFs) Readdir(path string) (files []string, err error) {
	// TODO can support param [options]
	_, ret := f.call("readdir", path)
	if hasError(ret) {
		return nil, JsErr(ret[0])
	}

	files = make([]string, ret[1].Length())
	for i := range files {
		files[i] = ret[1].Index(i).String()
	}
	return files, nil
}

func (f JsFs) Close(fd JsFd, callback func()) {
	args := []any{fd}
	if callback != nil {
		args = append(args, js.FuncOf(func(_ js.Value, _ []js.Value) any {
			callback()
			return nil
		}))
	}

	js.Value(Fs).Call("close", args...)
}

func (f JsFs) Open(path, flags string) (fd JsFd, err error) {
	_, ret := f.call("open", path, flags)
	if hasError(ret) {
		return JsFd{}, JsErr(ret[0])
	}

	return JsFd(ret[1]), nil
}

func (f JsFs) Utimes(path string, atime, mtime JsDate) (err error) {
	_, ret := f.call("utimes", path, atime, mtime)
	return wrapError(ret)
}

func (f JsFs) Futimes(fd JsFd, atime, mtime JsDate) (err error) {
	_, ret := f.call("futimes", fd, atime, mtime)
	return wrapError(ret)
}

func (f JsFs) Chown(path string, uid, gid uint64) (err error) {
	_, ret := f.call("chown", path, uid, gid)
	return wrapError(ret)
}

func (f JsFs) Fchown(fd JsFd, uid, gid uint64) (err error) {
	_, ret := f.call("fchown", fd, uid, gid)
	return wrapError(ret)
}

func (f JsFs) Chmod(path string, mode uint64) (err error) {
	_, ret := f.call("chmod", path, mode)
	return wrapError(ret)
}

func (f JsFs) Fchmod(fd JsFd, mode uint64) (err error) {
	_, ret := f.call("fchmod", fd, mode)
	return wrapError(ret)
}

func (f JsFs) Fsync(fd JsFd) (err error) {
	_, ret := f.call("fsync", fd)
	return wrapError(ret)
}

func (f JsFs) Write(fd JsFd, buffer []byte, offset, length, position int) (n int, err error) {
	buf := AllocBuffer(length)
	buf.Write(buffer, offset, length, 0)

	_, ret := f.call("write", fd, buf, 0, length, position)
	if hasError(ret) {
		return 0, JsErr(ret[0])
	}

	return ret[1].Int(), nil
}

func (f JsFs) Read(fd JsFd, buffer []byte, offset, length, position int) (n int, err error) {
	buf := AllocBuffer(length)
	_, ret := f.call("read", fd, buf, 0, length, position)
	if hasError(ret) {
		return 0, JsErr(ret[0])
	}

	n = ret[1].Int()
	buf.Read(buffer, offset, n, 0)

	return n, nil
}

func (f JsFs) ReadFile(path string) (bytes []byte, err error) {
	// TODO support options
	_, ret := f.call("readFile", path)
	if hasError(ret) {
		return nil, JsErr(ret[0])
	}

	return hex.DecodeString(ret[1].String())
}

func (f JsFs) WriteFile(filename string, data []byte) (err error) {
	// TODO support options
	_, ret := f.call("writeFile", filename, WrapBuffer(data))
	return wrapError(ret)
}

func (f JsFs) AppendFile(path string, p []byte) (err error) {
	// TODO support options
	_, ret := f.call("appendFile", path, WrapBuffer(p))
	return wrapError(ret)
}

func (f JsFs) Setxattr(path, name, value string) (err error) {
	// TODO support flag and complex valuue
	_, ret := f.call("setxattr", path, name, value)
	return wrapError(ret)
}

func (f JsFs) Fsetxattr(fd JsFd, name, value string) (err error) {
	// TODO support flag and complex valuue
	_, ret := f.call("fsetxattr", fd, name, value)
	return wrapError(ret)
}

func (f JsFs) Getxattr(path, name string) (value string, err error) {
	_, ret := f.call("getxattr", path, name, value)
	if hasError(ret) {
		return "", JsErr(ret[0])
	}

	return ret[1].String(), nil
}

func (f JsFs) Fgetxattr(fd JsFd, name string) (value string, err error) {
	_, ret := f.call("fgetxattr", fd, name)
	if hasError(ret) {
		return "", JsErr(ret[0])
	}

	return ret[1].String(), nil
}

func (f JsFs) Removexattr(path, name string) (err error) {
	_, ret := f.call("removexattr", path, name)
	return wrapError(ret)
}

func (f JsFs) Fremovexattr(fd JsFd, name string) (err error) {
	_, ret := f.call("fremovexattr", fd, name)
	return wrapError(ret)
}

func (f JsFs) Watch(filename string, options *WatchOption, listener func(event WatchEvent, filename string)) {
	args := []any{filename}
	if options != nil {
		args = append(args, options)
	}
	if listener != nil {
		args = append(args, js.FuncOf(func(this js.Value, args []js.Value) any {
			listener(WatchEvent(args[0].Int()), args[1].String())
			return nil
		}))
	}

	js.Value(Fs).Call("watch", args...)
}
