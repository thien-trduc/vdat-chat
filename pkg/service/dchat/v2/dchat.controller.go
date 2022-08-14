package dchat

import (
	_ "bytes"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"log"
	"net/http"
)

// Chat WebSocket godoc
// @Summary Chat websocket
// @Description chat group by websocket
// @Tags dchat
// @Param socketId path string true "socketId to know client"
// @Param token query string true "token to be join chat"
// @Accept  json
// @Produce  json
// @Router /message/{socketId} [get]
func ChatHandlr(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)
	// authenticate
	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
		return
	}

	param := mux.Vars(r)
	socketId := param["socketId"]
	if len(socketId) <= 0 {
		log.Println("Url Param 'socketId' is missing")
		return
	}

	v := r.URL.Query()
	paramOwner := v.Get("token")
	if len(paramOwner) <= 0 {
		log.Println("Url Param 'token' is missing")
		return
	}
	owner := auth.JWTparseOwner2(paramOwner)

	client := &Client{UserId: owner, SocketId: socketId, Broker: Wsbroker, Conn: conn, Send: make(chan []byte, 256)}
	client.Broker.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.

	go client.WritePump()
	go client.ReadPump()

}
