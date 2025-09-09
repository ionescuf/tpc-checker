package main

import (
	"go2-gkes-tpc/src/config"
	"go2-gkes-tpc/src/server"
)

func main() {
	c, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	server.Main(c)
}
