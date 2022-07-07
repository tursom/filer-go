package filer

import "syscall/js"

type (
	FileStats struct {
		Node    string // internal node id (unique)
		Dev     string // file system name
		Name    string // the entry's name (basename)
		Size    int    // file size in bytes
		Nlinks  int    // number of links
		Atime   JsDate // last access time as JS Date Object
		Mtime   JsDate // last modified time as JS Date Object
		Ctime   JsDate // creation time as JS Date Object
		AtimeMs uint64 // last access time as Unix Timestamp
		MtimeMs uint64 // last modified time as Unix Timestamp
		CtimeMs uint64 // creation time as Unix Timestamp
		Type    string // file type (FILE, DIRECTORY, SYMLINK),
		Gid     uint64 // group name
		Uid     uint64 // owner name
		Mode    uint64 // permissions
		Version uint64 // version of the node
	}
)

func (j FileStats) JsValue() js.Value {
	valueMap := make(map[string]any)
	valueMap["node"] = j.Node
	valueMap["dev"] = j.Dev
	valueMap["name"] = j.Name
	valueMap["size"] = j.Size
	valueMap["nlinks"] = j.Nlinks
	valueMap["atime"] = j.Atime
	valueMap["mtime"] = j.Mtime
	valueMap["ctime"] = j.Ctime
	valueMap["atimeMs"] = j.AtimeMs
	valueMap["mtimeMs"] = j.MtimeMs
	valueMap["ctimeMs"] = j.CtimeMs
	valueMap["type"] = j.Type
	valueMap["gid"] = j.Gid
	valueMap["uid"] = j.Uid
	valueMap["mode"] = j.Mode
	valueMap["version"] = j.Version
	return js.ValueOf(valueMap)
}

func GetFileStats(value js.Value) *FileStats {
	if value.Type() != js.TypeObject {
		return nil
	}

	return &FileStats{
		Node:    value.Get("node").String(),
		Dev:     value.Get("dev").String(),
		Name:    value.Get("node").String(),
		Size:    value.Get("size").Int(),
		Nlinks:  value.Get("nlinks").Int(),
		Atime:   JsDate(value.Get("atime")),
		Mtime:   JsDate(value.Get("mtime")),
		Ctime:   JsDate(value.Get("ctime")),
		AtimeMs: uint64(value.Get("atimeMs").Int()),
		MtimeMs: uint64(value.Get("mtimeMs").Int()),
		CtimeMs: uint64(value.Get("ctimeMs").Int()),
		Type:    value.Get("type").String(),
		Gid:     uint64(value.Get("gid").Int()),
		Uid:     uint64(value.Get("uid").Int()),
		Mode:    uint64(value.Get("mode").Int()),
		Version: uint64(value.Get("version").Int()),
	}
}
