package main

import (
	"log"

	"github.com/cbi-sh/messages/internal/app/router"
	"github.com/cbi-sh/messages/internal/app/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	r := router.New()

	// server.NewSMPPConnector("192.168.0.2:3736", "user", "password")
	s := server.New(":8027", r)
	s.Start()
}
