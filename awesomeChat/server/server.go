package main

import (
	fetch "src/awesomeChat/server/struct"
)

func main() {
	/*
		起飞
	*/

	server := new(fetch.Server)
	server.IP = [4]byte{127, 0, 0, 1}
	server.PORT = 8000
	server.Start()
}
