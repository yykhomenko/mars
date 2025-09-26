package main

import (
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
	http.Start()
}
