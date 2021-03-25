package main

import (
	"shiva/net"
)

func main() {
	s := net.NewServer("v0.2")

	s.Serve()
}
