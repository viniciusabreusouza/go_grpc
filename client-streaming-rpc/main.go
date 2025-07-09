package main

import (
	"example.com/m/client"
	"example.com/m/server"
)

func main() {
	go server.Run()
	client.Run()
}
