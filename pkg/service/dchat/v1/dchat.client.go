package dchat

import (
	"bytes"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	Broker *Broker

	// user ID, this will be parse from Access Token in production
	UserId string

	// socket id
	SocketId string

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of Outbound messages.
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Broker.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			sentry.CaptureException(err)
			log.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)

			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		var messageJSON Message
		_ = json.Unmarshal(message, &messageJSON)
		messageJSON.Client = c.UserId
		c.Broker.Inbound <- messageJSON
	}
}

// WritePump pumps messages from the Broker to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
//func (c *Client) WritePump() {
//	ticker := time.NewTicker(pingPeriod)
//	defer func() {
//		ticker.Stop()
//		c.Conn.Close()
//	}()
//	for {
//		select {
//		case message, ok := <-c.Send:
//			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if !ok {
//				// The Broker closed the channel.
//				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//
//			w, err := c.Conn.NextWriter(websocket.TextMessage)
//			if err != nil {
//				return
//			}
//			w.Write(message)
//
//			// Add queued chat messages to the current websocket message.
//			n := len(c.Send)
//			for i := 0; i < n; i++ {
//				w.Write(newline)
//				w.Write(<-c.Send)
//			}
//
//			if err := w.Close(); err != nil {
//				return
//			}
//		case <-ticker.C:
//			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//				return
//			}
//		}
//	}
//}
// Write writes a message with the given message type and payload.
func (c *Client) Write(mt int, payload []byte) error {
	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(mt, payload)
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The hub closed the channel.
				c.Write(websocket.CloseMessage, []byte{})
				return
			}

			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
