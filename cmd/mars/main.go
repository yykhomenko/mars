package main

import (
	"github.com/yykhomenko/mars/pkg/mars/api/http"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
)

func main() {

	//conf := config.NewConfig()

	router := router.NewRouter()
	//
	//smpp := smpp.NewSMPPConnector("localhost:3736", "user", "password", router)
	//smpp.Start()

	http := http.NewHTTPServer(":8080", router)
	http.Start()
}
