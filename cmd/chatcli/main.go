package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

func main() {
	//origin := "http://localhost/"
	//url := "ws://localhost:5000/ws"
	//ws, err := websocket.Dial(url, "", origin)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Print("message:")
	//messageReader := bufio.NewReader(os.Stdin)
	//messageInput, _ := messageReader.ReadString('\n')
	//messageInput = messageInput[:len(messageInput)-1]

	//messagePayload := model.MessagePayload{
	//
	//	Message:    "daioshdiophiqoh",
	//	SenderID:   "cuong",
	//	ReceiverID: "thien",
	//}
	//
	//if _, err := ws.Write(utils.ResponseWithByte(messagePayload)); err != nil {
	//	log.Fatal(err)
	//}
	//for {
	//	origin := "http://localhost/"
	//	url := "ws://localhost:5000/echo"
	//	ws, err := websocket.Dial(url, "", origin)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Print("message:")
	//	//messageReader := bufio.NewReader(os.Stdin)
	//	//messageInput, _ := messageReader.ReadString('\n')
	//	//messageInput = messageInput[:len(messageInput)-1]
	//
	//	messagePayload := model.MessagePayload{
	//
	//		Message:    "daioshdiophiqoh",
	//		SenderID:   "cuong",
	//		ReceiverID: "thien",
	//	}
	//
	//	if _, err := ws.Write(utils.ResponseWithByte(messagePayload)); err != nil {
	//		log.Fatal(err)
	//	}
	//}

	token, err := login()
	if err != nil {
		panic(err)
	}

	fmt.Println(token.AccessToken)

	//fmt.Println("Nhap nguoi gui : ")
	//senderReader := bufio.NewReader(os.Stdin)
	//senderInput, _ := senderReader.ReadString('\n')
	//senderInput = senderInput[:len(senderInput)-1]
	//
	//fmt.Println("Nhap nguoi nhan : ")
	//receivererReader := bufio.NewReader(os.Stdin)
	//receivererInput, _ := receivererReader.ReadString('\n')
	//receivererInput = receivererInput[:len(receivererInput)-1]

	//var serverURL = "ws://f8b423a4e609.ngrok.io/test"
	//var serverURL = "ws://localhost:5000/test"
	var serverURL = "ws://localhost:5000/user-online"
	if u := os.Getenv("SERVER_URL"); u != "" {
		serverURL = u
	}

	header := map[string][]string{"Authorization": {"Bearer " + token.AccessToken}}
	c, _, err := websocket.DefaultDialer.Dial(serverURL, header)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	//go func() {
	//	for {
	//		_, message, err := c.ReadMessage()
	//		if err != nil {
	//			log.Println("read:", err)
	//			return
	//		}
	//		log.Printf("recv: %s", message)
	//		time.Sleep(2 * time.Microsecond)
	//	}
	//}()
	//for {
	//	fmt.Println("message:")
	//	messageReader := bufio.NewReader(os.Stdin)
	//	messageInput, _ := messageReader.ReadString('\n')
	//	messageInput = messageInput[:len(messageInput)-1]
	//	messagePayload := model.MessagePayload{
	//
	//		Message:    messageInput,
	//		SenderID:   senderInput,
	//		ReceiverID: receivererInput,
	//	}
	//	if err = c.WriteJSON(messagePayload); err != nil {
	//		log.Fatal(err)
	//	}
	//}

}
