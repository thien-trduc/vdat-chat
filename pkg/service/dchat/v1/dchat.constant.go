package dchat

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	// Time allowed to Write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// event
const (
	SUBCRIBE      = "subcribe_group"
	SEND          = "send_text"
	LOADOLDMESS   = "load_old_mess"
	RELY          = "rely_message"
	LOADCHILDMESS = "load_child_mess"
	DELETE        = "delete_mess"
	FILE          = "FILE"
	TEXT          = "TEXT"
	IMAGE         = "IMAGE"
	UPDATE        = "update_mess"
)

var (
	newline = []byte{'\n'}

	space = []byte{' '}
)

var WsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
