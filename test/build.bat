set GOOS=js
set GOARCH=wasm
go build -trimpath -v -trimpath -o ./html/test.wasm .
