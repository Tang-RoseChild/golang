package main

import (
	"fmt"
	"redistest/chat/server"
)

func main() {
	err := server.Run(":23456")
	if err != nil {
		fmt.Println("err : ", err.Error())
	}
}
