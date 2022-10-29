package main

import (
	conn "src/awesomeChat/user/struct"
)

func main() {
	conn := new(conn.Conn)
	conn.IP = [4]byte{127, 0, 0, 1}
	conn.PORT = 8000
	conn.UID = 0
	conn.Locker = true
	conn.Start()
}
