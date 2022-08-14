package service

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"log"
//	"net/http"
//	"strings"
//	"time"
//
//	"github.com/dgrijalva/jwt-go"
//
//	"github.com/gorilla/websocket"
//)
//
//const (
//	// Time allowed to write a message to the peer.
//	writeWait = 10 * time.Second
//
//	// Time allowed to read the next pong message from the peer.
//	pongWait = 60 * time.Second
//
//	// Send pings to peer with this period. Must be less than pongWait.
//	pingPeriod = (pongWait * 9) / 10
//
//	// Maximum message size allowed from peer.
//	maxMessageSize = 512
//)
//
//var (
//	newline = []byte{'\n'}
//	space   = []byte{' '}
//)
//
//var broker *TestBroker
//
//var upgrader = websocket.Upgrader{
//	ReadBufferSize:  1024,
//	WriteBufferSize: 1024,
//	CheckOrigin: func(r *http.Request) bool {
//		return true
//	},
//}
//
//type Message struct {
//	To     string `json:"to"`
//	From   string `json:"from"`
//	Body   string `json:"body"`
//	status string
//}
//
//type TestBroker struct {
//	// Registered clients.
//	clients map[*Client]bool
//
//	// Inbound message from the clients.
//	inbound chan Message
//
//	// Outbound message that need to send to clients.
//	outbound chan Message
//
//	// Register request from the clients.
//	register chan *Client
//
//	// Unregister request from clients.
//	unregister chan *Client
//
//	messageRepository []*Message
//}
//
//func newbroker() *TestBroker {
//	return &TestBroker{
//		inbound:    make(chan Message),
//		outbound:   make(chan Message),
//		register:   make(chan *Client),
//		unregister: make(chan *Client),
//		clients:    make(map[*Client]bool),
//	}
//}
//
//func (b *TestBroker) run() {
//	// polling "new" message from repository
//	// and send to outbound channel to send to clients
//	// finally, marked message that sent to outbound channel as "done"
//	go func() {
//		for {
//			for idx, m := range b.messageRepository {
//				if m.status != "done" {
//					select {
//					case b.outbound <- *m:
//					default:
//						close(b.outbound)
//					}
//
//					b.messageRepository[idx].status = "done"
//				}
//			}
//
//			time.Sleep(200 * time.Millisecond)
//		}
//	}()
//
//	for {
//		select {
//		case client := <-b.register:
//			b.clients[client] = true
//		case client := <-b.unregister:
//			if _, ok := b.clients[client]; ok {
//				delete(b.clients, client)
//				close(client.send)
//			}
//		case message := <-b.inbound:
//
//			b.messageRepository = append(b.messageRepository, &message)
//			fmt.Printf("%+v, %d\n", message, len(b.messageRepository))
//
//		case message := <-b.outbound:
//			fmt.Println("send")
//			for client := range b.clients {
//				if client.userID == message.To {
//					msg, _ := json.Marshal(message)
//					select {
//					case client.send <- msg:
//					default:
//						close(client.send)
//						delete(b.clients, client)
//					}
//				}
//			}
//		}
//
//	}
//}
//
//type Client struct {
//	broker *TestBroker
//
//	// user ID, this will be parse from Access Token in production
//	userID string
//
//	// The websocket connection.
//	conn *websocket.Conn
//
//	// Buffered channel of outbound message.
//	send chan []byte
//}
//
//func (c *Client) readPump() {
//	defer func() {
//		c.broker.unregister <- c
//		c.conn.Close()
//	}()
//	c.conn.SetReadLimit(maxMessageSize)
//	c.conn.SetReadDeadline(time.Now().Add(pongWait))
//	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
//	for {
//		_, message, err := c.conn.ReadMessage()
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				log.Printf("error: %v", err)
//			}
//			break
//		}
//		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
//
//		var messageJSON Message
//		_ = json.Unmarshal(message, &messageJSON)
//		messageJSON.From = c.userID
//
//		c.broker.inbound <- messageJSON
//	}
//}
//
//// writePump pumps message from the broker to the websocket connection.
////
//// A goroutine running writePump is started for each connection. The
//// application ensures that there is at most one writer to a connection by
//// executing all writes from this goroutine.
//func (c *Client) writePump() {
//	ticker := time.NewTicker(pingPeriod)
//	defer func() {
//		ticker.Stop()
//		c.conn.Close()
//	}()
//	for {
//		select {
//		case message, ok := <-c.send:
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if !ok {
//				// The broker closed the channel.
//				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
//				return
//			}
//
//			w, err := c.conn.NextWriter(websocket.TextMessage)
//			if err != nil {
//				return
//			}
//			w.Write(message)
//
//			// Add queued chat message to the current websocket message.
//			n := len(c.send)
//			for i := 0; i < n; i++ {
//				w.Write(newline)
//				w.Write(<-c.send)
//			}
//
//			if err := w.Close(); err != nil {
//				return
//			}
//		case <-ticker.C:
//			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
//			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
//				return
//			}
//		}
//	}
//}
//
//func TestHandler(w http.ResponseWriter, r *http.Request) {
//	// authenticate
//	var (
//		claims jwt.StandardClaims
//
//		// https://accounts.vdatlab.com/auth/realms/vdatlab.com/protocol/openid-connect/certs
//		publicKey = "MIICpTCCAY0CBgFrPLdvYjANBgkqhkiG9w0BAQsFADAWMRQwEgYDVQQDDAt2ZGF0bGFiLmNvbTAeFw0xOTA2MDkxNDQ4MDNaFw0yOTA2MDkxNDQ5NDNaMBYxFDASBgNVBAMMC3ZkYXRsYWIuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAleIkHoO6Q0GRQ4POIAKmN5Ev3zfAm8raTJQ1e/CbTXW4FQ0kDS9YPhLXPwcdnbxiL3rSGgz7+iWcq/Ix7yExuNbSyqDUjLUJSU6I9JvB1YP8GSaO8d996+TVCDC8E/VSID6wmfWbMNb5Ns6Y7YY/HAhj9zc73ObErvi0NV0BjeYAVOBqJKKgl9cHfyBshr+kpC/7nrbTRnAP7JQhKrQF6wBTKQiuJlEyYqvi1ugCRBYg2BZLPtTry+Kineb1DT8ynmxJjKMtr9hU0dsLPJpqW/4DWwNOarLOBP/K9WkfR2LUxbrm41goSTjJbz6s7f/Mvn/gDLjGjIsdlFP3Y7I2lwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQBXU5Awwhv/cYJKCdSUzmtXpXty8KrdrHaDNa8potDXlEc2JrK3wHyFRwwfpBhkaicP0LllxRHUGUNsWFnggae1fudc75fysZ16NPH7VJlUuyV96K06K4v1aM5VCSWl5djky7rtyfi2W9iH2ddWZvCeSyFsSgCD4P5GjgYpsLy27g/cvdJJAdp/b7bweVDI1grlBtnInxLUPhJ4cnoNw3crh7twqKgG6F3GmZc2Hjl45LdlxBFfftDUYH66D1X0mdoipQCbg4JWlIxUZHVjJDIrSIlwnRMwjzCm7MUYv0ySmvsxgoNVI2NuFU6A/F7zlyVkDkmO4ilp4BueRtBKb7yR"
//	)
//	{
//		var tok string
//		queryToken := r.URL.Query()["token"]
//		if len(queryToken) > 0 {
//			tok = queryToken[0]
//		}
//
//		var accessToken string
//		h := r.Header.Get("Authorization")
//		a := strings.Split(h, " ")
//		if len(a) == 2 {
//			accessToken = a[1]
//		} else if tok != "" {
//			accessToken = tok
//		} else {
//			w.WriteHeader(400)
//			return
//		}
//
//		token, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (interface{}, error) {
//			return jwt.ParseRSAPublicKeyFromPEM([]byte("-----BEGIN CERTIFICATE-----\n" + publicKey + "\n-----END CERTIFICATE-----"))
//		})
//
//		if err != nil {
//			fmt.Println("CANNOT parse token: ", err)
//			w.WriteHeader(401)
//			return
//		}
//
//		if token.Valid {
//			fmt.Println("Token is VALID")
//		} else if ve, ok := err.(*jwt.ValidationError); ok {
//			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
//				fmt.Println("That's not even a token")
//			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
//				// Token is either expired or not active yet
//				fmt.Println("Timing is everything")
//			} else {
//				fmt.Println("Couldn't handle this token:", err)
//			}
//		} else {
//			fmt.Println("Couldn't handle this token:", err)
//		}
//	}
//
//	if broker == nil {
//		broker = &TestBroker{
//			inbound:    make(chan Message),
//			outbound:   make(chan Message),
//			register:   make(chan *Client),
//			unregister: make(chan *Client),
//			clients:    make(map[*Client]bool),
//		}
//
//		go broker.run()
//	}
//
//	conn, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	fmt.Println("subject: " + claims.Subject)
//
//	client := &Client{userID: claims.Subject, broker: broker, conn: conn, send: make(chan []byte, 256)}
//	client.broker.register <- client
//
//	// Allow collection of memory referenced by the caller by doing all work in
//	// new goroutines.
//	go client.writePump()
//	go client.readPump()
//}
