package main

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/pogo.proto

import (
	"github.com/mt-inside/pogo/pogo/cmd"
)

func main() {
	cmd.Execute()
}
