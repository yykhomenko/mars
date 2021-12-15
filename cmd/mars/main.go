package main

import (
	"log"

	"github.com/yykhomenko/mars/internal/api/http"
	"github.com/yykhomenko/mars/internal/service/router"
	"github.com/yykhomenko/mars/internal/service/smpp"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	r := router.NewRouter()

	smpp := smpp.NewSMPPConnector("localhost:3736", "user", "password", r)
	smpp.Start()

	http := http.NewHTTPServer(":8080", r)
	http.Start()
}
