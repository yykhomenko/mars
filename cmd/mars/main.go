package main

import (
	"fmt"
	"time"

	"github.com/yykhomenko/mars/pkg/mars/api/http"
	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/service/hash"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
	"github.com/yykhomenko/mars/pkg/mars/service/sis"
)

func main() {

	config := config.NewConfig()
	//fmt.Println(config)

	router := router.NewRouter(config)
	//
	//smpp := smpp.NewSMPPConnector("localhost:3736", "user", "password", router)
	//smpp.Start()

	hashConnector := hash.NewHashConnector(config)
	sisConnector := sis.NewSisConnector(config)

	http := http.NewHTTPServer(config, hashConnector, sisConnector, router)

	prevNum := router.GetNum()
	go func() {
		for range time.Tick(1 * time.Second) {
			currentNum := router.GetNum()
			if prevNum != currentNum {
				fmt.Printf("%d tps\n", currentNum-prevNum)
			}

			//fmt.Println(sisConnector.GetSubscriber("380670000001"))

			prevNum = currentNum
		}
	}()

	http.Start()
}
