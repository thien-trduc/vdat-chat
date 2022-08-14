package dchat

import (
	"bytes"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/websocket"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/cors"
	"log"
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
	maxMessageSize = 2048
)

var (
	newline    = []byte{'\n'}
	space      = []byte{' '}
	WsUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	UserId   string
	ClientId string
	Send     chan []byte
	Identify bool
	Conn     *websocket.Conn
	Broker   *Broker
}

func (c *Client) ReadPump() {
	defer func() {
		c.Broker.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	err := c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
		return
	}
	c.Conn.SetPongHandler(func(string) (err error) {
		err = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			return
		}
		return
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				sentry.CaptureException(err)
				log.Printf("Exception : %s", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var messageJSON MessageRequest
		_ = json.Unmarshal(message, &messageJSON)
		if messageJSON.RequestType == SUBSCRIBE {
			c.ClientId = messageJSON.ClientId
			c.Broker.Inbound <- messageJSON
		} else {
			if c.Identify == true {
				c.Broker.Inbound <- messageJSON
			}
		}
	}
}

// TODO : Write writes a message with the given message type and payload.
func (c *Client) Write(mt int, payload []byte) (err error) {
	err = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return
	}
	err = c.Conn.WriteMessage(mt, payload)
	if err != nil {
		return
	}
	return
}

// TODO :  WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if c.Identify == true {
				if !ok {
					// The hub closed the channel.
					err := c.Write(websocket.CloseMessage, []byte{})
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
					}
					return
				}

				err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
				w, err := c.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
				_, err = w.Write(message)
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
				// Add queued chat messages to the current websocket message.
				n := len(c.Send)
				for i := 0; i < n; i++ {
					_, err = w.Write(newline)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
						return
					}
					_, err = w.Write(<-c.Send)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
						return
					}
				}

				if err := w.Close(); err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
			} else {
				if !ok {
					// The hub closed the channel.
					err := c.Write(websocket.CloseMessage, []byte{})
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
					}
					return
				}

				err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
				w, err := c.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
				_, err = w.Write(message)
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
				// Add queued chat messages to the current websocket message.
				n := len(c.Send)
				for i := 0; i < n; i++ {
					_, err = w.Write(newline)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
						return
					}
					_, err = w.Write(<-c.Send)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
						return
					}
				}

				if err := w.Close(); err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
					return
				}
			}
		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				sentry.CaptureException(err)
				log.Printf("Exception : %s", err)
				return
			}
		}
	}
}

func ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	cors.SetupResponse(&w, r)

	conn, err := WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		sentry.CaptureException(err)
		log.Println(err)
		return
	}

	client := &Client{Broker: Wsbroker, Conn: conn, Identify: false, Send: make(chan []byte, 256)}
	client.Broker.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
