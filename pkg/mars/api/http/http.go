package http

import (
	"net/http"
	"time"

	"github.com/yykhomenko/mars/pkg/mars/config"
	"github.com/yykhomenko/mars/pkg/mars/entity"
	"github.com/yykhomenko/mars/pkg/mars/service/hash"
	"github.com/yykhomenko/mars/pkg/mars/service/router"
)

type HTTPServer struct {
	conf   *config.Config
	router router.Router
}

func NewHTTPServer(conf *config.Config, hashConnector *hash.HashConnector, router router.Router) *HTTPServer {
	s := &HTTPServer{
		conf:   conf,
		router: router,
	}

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		r.ParseForm()
		from := r.FormValue("from")
		to := r.FormValue("to")
		text := r.FormValue("text")

		message := &entity.Message{
			From: from,
			To:   to,
			Text: text,
		}

		conf.Log.Println("1", message)

		hash, err := hashConnector.GetHash(message.To)
		if err != nil {
			return
		}

		conf.Log.Println("2", hash)

		s.router.Route(message)

		conf.Log.Printf("http: rx: duration: %s", time.Since(start))
	})

	return s
}

func (s *HTTPServer) Start() error {
	s.conf.Log.Println("HTTP server listen:", s.conf.MarsApiAddr)
	return http.ListenAndServe(s.conf.MarsApiAddr, nil)
}
