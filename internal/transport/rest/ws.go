package rest

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkAllowedOrigins,
}

// checkAllowedOrigins returns true if the origin is allowed.
func checkAllowedOrigins(r *http.Request) bool {
	allowedOrigins := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
	}
	allowedOrigins = append(allowedOrigins, os.Getenv("FRONTEND_URL"))

	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	for _, allowedOrigin := range allowedOrigins {
		if strings.EqualFold(allowedOrigin, u.Scheme+"://"+u.Host) {
			return true
		}
	}

	return false
}

func upgradeToWs(writer http.ResponseWriter, request *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
