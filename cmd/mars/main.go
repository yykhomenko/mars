package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yykhomenko/mars/pkg/mars/api/http"
	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/service/hash"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
	"github.com/yykhomenko/mars/pkg/mars/service/sis"
)

type application struct {
	config        *config.Config
	log           *logrus.Logger
	hashConnector *hash.HashConnector
	sisConnector  *sis.SisConnector
	router        *router.Router
	httpServer    *http.HTTPServer
}

func main() {

	config := config.NewConfig()
	hashConnector := hash.NewHashConnector(config)
	sisConnector := sis.NewSisConnector(config)
	router := router.NewRouter(config)
	http := http.NewHTTPServer(config, hashConnector, sisConnector, router)

	prevNum := router.GetNum()
	go func() {
		for range time.Tick(1 * time.Second) {
			currentNum := router.GetNum()
			if prevNum != currentNum {
				fmt.Printf("%d tps\n", currentNum-prevNum)
			}
			prevNum = currentNum
		}
	}()

	http.Start()
}

//
//smpp := smpp.NewSMPPConnector("localhost:3736", "user", "password", router)
//smpp.Start()
