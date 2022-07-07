package main

import (
	"encoding/json"
	"fmt"

	"filer-go/filer"
)

func log(name string, obj any) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("marshal require fail, err: ", err)
		return
	}

	fmt.Println(fmt.Sprintf("name: %s, obj: %v, json: %s", name, obj, string(bytes)))
}

func main() {
	fmt.Println("hello, world!")
	fd, err := filer.Fs.Open("/test", "w+")
	if err != nil {
		fmt.Println("open file failed, err: ", err)
		return
	}

	//n, err := filer.Fs.Write(fd, []byte("hello, world!"), 0, 13, 0)
	//if err != nil {
	//	fmt.Println("write file failed, err: ", err)
	//	return
	//}
	//fmt.Println("write file ", n, " bytes")

	buf := make([]byte, 100)
	n, err := filer.Fs.Read(fd, buf, 0, 100, 0)
	if err != nil {
		fmt.Println("read file failed, err: ", err)
		return
	}
	fmt.Println("read file ", n, " bytes")
	fmt.Println("read content: ", string(buf[:n]))
}
