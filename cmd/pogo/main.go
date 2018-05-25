package main

//go:generate protoc -I ../../api --go_out=plugins=grpc:../../api ../../api/pogo.proto

import (
	"github.com/mt-inside/pogo/cmd/pogo/commands"
)

func main() {
	commands.Execute()
}
