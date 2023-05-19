package proxy

import (
	"net/http"

	"github.com/hashicorp/yamux"
	"golang.org/x/net/websocket"
)

func Gateway(gateway string, h http.Handler) {
	config, err := websocket.NewConfig(gateway, gateway)
	if err != nil {
		panic(err)
	}
	conn, err := websocket.DialConfig(config)
	if err != nil {
		panic(err)
	}

	server, err := yamux.Server(conn, nil)
	if err != nil {
		panic(err)
	}

	err = http.Serve(server, h)
	if err != nil {
		panic(err)
	}
}
