package main

import (
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", "lo1.zzu.wiki:8080")
	if err != nil {
		log.Fatalln(err)
	}
	_ = lis.Close()
}

type Init struct {
	DSN string
}
