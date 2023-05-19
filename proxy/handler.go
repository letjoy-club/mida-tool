package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/hashicorp/yamux"
	"golang.org/x/net/websocket"
)

var reverseProxy *httputil.ReverseProxy

var userIDs = []string{}

func RedirectProxy(uid string, w http.ResponseWriter, r *http.Request) bool {
	p := reverseProxy
	if p == nil {
		return false
	}

	for _, id := range userIDs {
		if id == uid {
			p.ServeHTTP(w, r)
			return true
		}
	}
	return false
}

func ProxyHandler() *websocket.Server {
	// 这个变量无关紧要
	u, _ := url.Parse("http://localhost")
	reverseProxy = httputil.NewSingleHostReverseProxy(u)

	wsServer := websocket.Server{
		Config: websocket.Config{},
		Handler: func(w *websocket.Conn) {
			defer w.Close()

			client, err := yamux.Client(w, nil)
			if err != nil {
				return
			}
			defer client.Close()

			query := w.Request().URL.Query()
			strs := query.Get("userIds")
			if strs == "" {
				return
			}

			userIDs = strings.Split(strs, ",")

			reverseProxy = httputil.NewSingleHostReverseProxy(u)
			reverseProxy.Transport = &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return client.Open()
				},
			}

			defer func() {
				reverseProxy = nil
			}()

			fmt.Println("accept reverse proxy connection")

			for {
				_, err := client.Ping()
				if err != nil {
					fmt.Println(err)
					return
				}
			}

		},
	}
	return &wsServer
}
