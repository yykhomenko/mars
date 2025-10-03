package main

import (
	"fmt"
	"time"

	"github.com/yykhomenko/mars/pkg/mars/api/http"
	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/service/hash"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
)

func main() {

	config := config.NewConfig()
	//fmt.Println(config)

	router := router.NewRouter(config)
	//
	//smpp := smpp.NewSMPPConnector("localhost:3736", "user", "password", router)
	//smpp.Start()

	hashConnector := hash.NewHashConnector(config)

	http := http.NewHTTPServer(config, hashConnector, router)

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
