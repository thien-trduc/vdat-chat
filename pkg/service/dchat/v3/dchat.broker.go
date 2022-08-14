package dchat

import (
	"context"
	"fmt"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/groups/v2"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/useronline/v2"
	"time"
)

type Broker struct {
	Clients    map[*Client]bool
	Inbound    chan MessageRequest
	Register   chan *Client
	Unregister chan *Client
	Ctx        context.Context
}

var Wsbroker = &Broker{
	Clients:    make(map[*Client]bool),
	Inbound:    make(chan MessageRequest),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Ctx:        context.Background(),
}

func (b *Broker) Run() {
	timeoutContext := time.Duration(2) * time.Second
	messRepo := message_service.NewRepoImpl(database.DB)
	messService := message_service.NewServiceImpl(messRepo, timeoutContext)
	userOnRepo := useronline.NewRepoImpl(database.DB)
	userOnService := useronline.NewServiceImpl(userOnRepo, timeoutContext)
	groupRepo := groups.NewRepoImpl(database.DB)
	userRepo := userdetail.NewRepoImpl(database.DB)
	userService := userdetail.NewServiceImpl(userRepo, timeoutContext)
	groupService := groups.NewServiceImpl(groupRepo, userService, timeoutContext, messService)
	events := NewEventImpl(b, b.Ctx, userOnService, messService, groupService)

	for {
		select {
		case client := <-b.Register:
			client.Identify = false
			b.Clients[client] = true
			fmt.Println("client " + client.UserId + " is connected")
		case client := <-b.Unregister:
			if _, ok := b.Clients[client]; ok {
				delete(b.Clients, client)
				close(client.Send)
			}
		case message := <-b.Inbound:
			switch message.RequestType {
			case SUBSCRIBE:
				events.SubscribeGroup(message)
			case LOAD_OLD_MESSAGE:
				events.LoadOldMessage(message)
			case SEND_MESSAGE:
				events.SendMessage(message)
			case UPDATE_MESSAGE:
				events.UpdateMessage(message)
			case DELETE_MESSAGE:
				events.DeleteMessage(message)
			case LOAD_CHILD_MESSAGE:
				events.LoadChildMessage(message)
			}
		}
	}
}
