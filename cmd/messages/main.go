package messages

import (
	"log"

	"github.com/cbi-sh/messages/internal/app/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	server.NewHttpServer(":8027")
	// server.NewSMPPConnector("192.168.0.2:3736", "user", "password")
}
