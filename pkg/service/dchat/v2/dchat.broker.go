package dchat

import (
	"context"
	"fmt"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/useronline/v2"
	"log"
	"time"
)

//var defalutParnetID = -1

type Broker struct {
	// 1 group va nhieu client connect toi
	Clients map[*Client]bool

	// Inbound messages from the Clients.
	Inbound chan Message

	// Outbound messages that need to Send to Clients.
	Outbound chan Message

	// Register request from the Clients.
	Register chan *Client

	// Unregister request from Clients.
	Unregister chan *Client

	MessageRepository []*Message

	Ctx context.Context
}

var Wsbroker = &Broker{
	Clients:    make(map[*Client]bool),
	Inbound:    make(chan Message),
	Outbound:   make(chan Message),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Ctx:        context.Background(),
}

// luong nay dc khoi tao khi chay truong trinh
// Broker la noi tiep nhan client khi client do mo connect
//nhan cac message dc gui len cua client va tra message dc gui di den client nhan
func (b *Broker) Run() {
	// polling "new" message from repository
	// and Send to Outbound channel to Send to Clients
	// finally, marked message that sent to Outbound channel as "done"
	timeoutContext := time.Duration(2) * time.Second
	messRepo := message_service.NewRepoImpl(database.DB)
	messService := message_service.NewServiceImpl(messRepo, timeoutContext)

	userOnRepo := useronline.NewRepoImpl(database.DB)
	userOnService := useronline.NewServiceImpl(userOnRepo, timeoutContext)

	events := NewEventImpl(b, b.Ctx, userOnService, messService)

	go func() {
		for {
			for idx, m := range b.MessageRepository {
				if m.Data.Status != "done" {
					log.Printf("idx : %d, mess : %v", idx, m)
					select {
					case b.Outbound <- *m:
						b.MessageRepository[idx].Data.Status = "done"
					default:
					}
				}
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()

	for {
		select {
		case client := <-b.Register:
			// khi client dang nhap thi broker se lay group dua tren idgroup cua client, neu chua co thi tao group vao broker
			b.Clients[client] = true
			fmt.Println("client " + client.UserId + " is connected")
		case client := <-b.Unregister:
			// khi client dang xuat khoi group, delete client khoi group
			if _, ok := b.Clients[client]; ok {
				//delete in database when client off
				delete(b.Clients, client)
				close(client.Send)
			}
			//khi co tin nhan dc gui di , message se duoc day vao inbound va day vao MessageRepository
		case message := <-b.Inbound:
			b.MessageRepository = append(b.MessageRepository, &message)
			//fmt.Printf("%+v, %d\n", message, len(b.MessageRepository))
		case message := <-b.Outbound:
			switch message.TypeEvent {
			case SEND:
				events.SendMessage(message)
			case RELY:
				events.ReplyMessage(message)
			case DELETE:
				events.DeleteMessage(message)
			case LOADCHILDMESS:
				events.LoadChildMessage(message)
			case SUBCRIBE:
				events.SubscribeGroup(message)
			case LOADOLDMESS:
				events.LoadOldMessage(message)
			case UPDATE:
				events.UpdateMessage(message)
			default:
			}
		}
	}
}
