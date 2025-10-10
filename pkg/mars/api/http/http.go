package http

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/entity"
	"github.com/yykhomenko/mars/pkg/mars/service/hash"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
	"github.com/yykhomenko/mars/pkg/mars/service/sis"
)

type HTTPServer struct {
	conf   *config.Config
	server *fiber.App
	router router.Router
}

func NewHTTPServer(conf *config.Config, hashConnector *hash.HashConnector, sisConnector *sis.SisConnector, router router.Router) *HTTPServer {

	rand.Seed(time.Now().Unix())
	app := fiber.New()

	s := &HTTPServer{
		conf:   conf,
		server: app,
		router: router,
	}

	app.Get("/messages", func(c *fiber.Ctx) error {

		//start := time.Now()

		from := c.FormValue("from")
		to := c.FormValue("to")
		text := c.FormValue("text")

		message := &entity.Message{
			From: from,
			To:   to,
			Text: text,
		}

		//conf.Log.Println("http: message: ", message)

		hashResp, err := hashConnector.GetHash(message.To)
		if err != nil {
			conf.Log.Warn("hash: response err:", err.Error())
		}
		////conf.Log.Println("hash: response: ", hashResp)
		//
		sisResp, err := sisConnector.GetSubscriber(message.To)
		if err != nil {
			conf.Log.Warn("sis: response err:", err.Error())
		}
		////conf.Log.Println("sis: response: ", sisResp)
		//
		message.SisType = sisResp.BillingType
		message.To = hashResp.Value

		s.router.Route(message)

		//conf.Log.Printf("http duration: %s", time.Since(start))

		return nil
	})

	return s
}

func (s *HTTPServer) Start() error {
	//s.conf.Log.Println("HTTP server listen:", s.conf.MarsApiAddr)
	//return http.ListenAndServe(s.conf.MarsApiAddr, nil)
	return s.server.Listen(s.conf.MarsApiAddr)
}
