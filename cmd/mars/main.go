package main

import (
	"log"

	"github.com/yykhomenko/mars/internal/app/router"
	"github.com/yykhomenko/mars/internal/app/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	r := router.New()

	smpp := server.NewSMPPConnector("192.168.0.2:3736", "user", "password", r)
	smpp.Start()

	http := server.NewHTTPConnector(":8080", r)
	http.Start()
}
