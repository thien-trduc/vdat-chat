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
		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				sentry.CaptureException(err)
				log.Printf("Exception : %s", err)
				return
			}
		}
	}
}
