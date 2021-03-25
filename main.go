package main

import (
	"shiva/net"
)

func main() {
	s := net.NewServer("v0.1")

	s.Serve()
}
