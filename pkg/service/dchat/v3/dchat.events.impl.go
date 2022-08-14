package dchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/auth"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/groups/v2"
	message_service "gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/upload/v1"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/useronline/v2"
	"golang.org/x/sync/errgroup"
	"log"
	"strconv"
)

type EventImpl struct {
	b             *Broker
	userOnService useronline.Service
	messService   message_service.Service
	groupService  groups.Service
	ctx           context.Context
}

func NewEventImpl(b *Broker, ctx context.Context, userOnService useronline.Service, messService message_service.Service, groupService groups.Service) Events {
	return &EventImpl{
		b:             b,
		userOnService: userOnService,
		messService:   messService,
		groupService:  groupService,
		ctx:           ctx,
	}
}

// * 1
// TODO : Xac Thuc User
func (e *EventImpl) SubscribeGroup(message MessageRequest) {
	var msg []byte
	response := MessageResponse{ResponseType: SUBSCRIBED}
	bodyResponse := SubscribeResponseBody{}

	bodyRequest := SubscribeMessageBody{}
	reqBodyBytes := new(bytes.Buffer)

	err := json.NewEncoder(reqBodyBytes).Encode(message.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	err = json.Unmarshal(reqBodyBytes.Bytes(), &bodyRequest)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	owner := auth.JWTparseOwner2(bodyRequest.AccessToken)
	for client := range e.b.Clients {
		log.Printf("Client : %v", client)
		if client.ClientId == message.ClientId && client.Identify == false {
			if len(owner) <= 0 {
				bodyResponse.Subscribed = false
				response.Body = bodyResponse
				msg, err = json.Marshal(response)
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s \n", err)
				}

				select {
				case client.Send <- msg:
					close(client.Send)
					delete(e.b.Clients, client)
				}
			} else {
				bodyResponse.Subscribed = true
				response.Body = bodyResponse
				client.Identify = true
				msg, err = json.Marshal(response)
				select {
				case client.Send <- msg:
					newClient := Client{
						UserId:   owner,
						ClientId: client.ClientId,
						Send:     client.Send,
						Identify: client.Identify,
						Conn:     client.Conn,
						Broker:   client.Broker,
					}
					newClientPoint := &newClient
					e.b.Clients[newClientPoint] = true
					delete(e.b.Clients, client)
				default:
					delete(e.b.Clients, client)
					close(client.Send)
				}
			}
			break
		}
	}

}

// * 2
// TODO : Load tin nhan cu
func (e *EventImpl) LoadOldMessage(message MessageRequest) {
	var msg []byte
	var parentId int

	response := MessageResponse{ResponseType: MESSAGE}
	bodyResponse := NewMessageResponse{}

	bodyRequest := LoadOldMessageBody{}
	reqBodyBytes := new(bytes.Buffer)

	err := json.NewEncoder(reqBodyBytes).Encode(message.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	err = json.Unmarshal(reqBodyBytes.Bytes(), &bodyRequest)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}

	oldMess, err := e.messService.LoadContinueMessageHistoryService(e.ctx, bodyRequest.LastMessageId, bodyRequest.GroupId)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}

	for client := range e.b.Clients {
		if client.UserId == message.SenderId && client.ClientId == message.ClientId && client.Identify == true {
			for _, o := range oldMess {
				if o.DeletedAt != nil {
					o.Content = ""
				} else {
					if o.Type == IMAGE_MESSAGE {
						bucketName := "group-" + strconv.Itoa(bodyRequest.GroupId)
						getLinkFile, err := upload.GetFileService(bucketName, o.Content)
						if err != nil {
							sentry.CaptureException(err)
							log.Printf("Exception : %s", err)
						}
						o.Content = getLinkFile
					}
				}

				if bodyRequest.ParentMessageId > -1 {
					parentId = o.ParentId
				} else {
					parentId = defalutParnetID
				}
				//userInfo := userdetail.GetUserById(o.SubjectSender)
				mess := Message{
					Id:                int(o.ID),
					ParentId:          parentId,
					GroupId:           o.IdGroup,
					Sender:            o.SubjectSender.ID,
					Message:           o.Content,
					MessageType:       o.Type,
					TotalChildMessage: o.NumChildMess,
					CreatedAt:         o.CreatedAt,
					UpdatedAt:         o.UpdatedAt,
					DeletedAt:         o.DeletedAt,
					UserInfo:          o.SubjectSender,
				}
				bodyResponse.Message = mess
				bodyResponse.GroupId = o.IdGroup
				response.Body = bodyResponse

				msg, err = json.Marshal(response)
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
				}
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(e.b.Clients, client)
				}
			}

		}
	}
}

// * 3
// TODO : Gui Tin Nhan
func (e *EventImpl) SendMessage(message MessageRequest) {
	var msg []byte

	response := MessageResponse{ResponseType: NEW_MESSAGE}
	bodyResponse := NewMessageResponse{}

	bodyRequest := SendMessageBody{}
	reqBodyBytes := new(bytes.Buffer)

	err := json.NewEncoder(reqBodyBytes).Encode(message.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	err = json.Unmarshal(reqBodyBytes.Bytes(), &bodyRequest)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}

	var newMess message_service.Dto
	newMessChan := make(chan message_service.Dto)
	g, ctx := errgroup.WithContext(e.ctx)
	g.Go(func() (err error) {
		payload := message_service.PayLoad{
			SubjectSender: message.SenderId,
			Content:       bodyRequest.Content,
			IdGroup:       bodyRequest.GroupId,
			Type:          bodyRequest.MessageType,
		}
		if bodyRequest.ParentMessageId > -1 {
			payload.ID = bodyRequest.ParentMessageId
			newMess, err = e.messService.AddRelyService(ctx, payload)
			if err != nil {
				return
			}
		} else {
			newMess, err = e.messService.AddMessageService(ctx, payload)
			if err != nil {
				return
			}
		}
		_, err = e.groupService.UpdateGroupWhenHaveAction(ctx, bodyRequest.GroupId)
		if err != nil {
			fmt.Println(err)
			sentry.CaptureException(err)
			return
		}
		newMessChan <- newMess
		return

	})
	// check type là hình
	if bodyRequest.MessageType == IMAGE_MESSAGE {
		fmt.Println(bodyRequest.Content)
		fmt.Println("đây là link 1")
		bucketName := "group-" + strconv.Itoa(bodyRequest.GroupId)
		getLinkFile, _ := upload.GetFileService(bucketName, bodyRequest.Content)
		bodyRequest.Content = getLinkFile
		fmt.Println(getLinkFile)
		fmt.Println("đây là link 2")
	}

	g.Go(func() (err error) {
		userOn, err := e.userOnService.GetListUSerOnlineByGroupService(ctx, bodyRequest.GroupId)
		if err != nil {
			return
		}
		select {
		case n := <-newMessChan:
			for client := range e.b.Clients {
				for _, u := range userOn {
					if u.UserID == client.UserId && u.SocketID == client.ClientId && client.Identify == true {
						fmt.Println("check")
						if bodyRequest.MessageType == CALL_MESSAGE {
							fmt.Println(bodyRequest.GroupId)
							bodyRequest.Content = bodyRequest.Content + "/call;meetingId=" + strconv.Itoa(bodyRequest.GroupId) + ";peerId=" + message.SenderId
							// use http://localhost:4200/call;meetingId=07927fc8-af0a-11ea-b338-064f26a5f90a;userId=alice;peerId=bob
							// and http://localhost:4200/call;meetingId=07927fc8-af0a-11ea-b338-064f26a5f90a;userId=bob;peerId=alice
						}
						if n.ParentId == int(n.ID) {
							n.ParentId = defalutParnetID
						}
						//userInfo := userdetail.GetUserById(n.SubjectSender)
						mess := Message{
							Id:                int(n.ID),
							ParentId:          n.ParentId,
							GroupId:           n.IdGroup,
							Sender:            n.SubjectSender.ID,
							Message:           bodyRequest.Content,
							MessageType:       n.Type,
							TotalChildMessage: n.NumChildMess,
							CreatedAt:         n.CreatedAt,
							UpdatedAt:         n.UpdatedAt,
							DeletedAt:         n.DeletedAt,
							UserInfo:          n.SubjectSender,
						}
						bodyResponse.Message = mess
						bodyResponse.GroupId = n.IdGroup
						response.Body = bodyResponse

						msg, err = json.Marshal(response)
						if err != nil {
							return err
						}
						select {
						case client.Send <- msg:
						default:
							close(client.Send)
							delete(e.b.Clients, client)
						}
					}
				}
			}
		}
		return
	})
	if err := g.Wait(); err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}
	close(newMessChan)
}

// * 4
func (e *EventImpl) UpdateMessage(message MessageRequest) {
	var msg []byte

	response := MessageResponse{ResponseType: UPDATED_MESSAGE}
	bodyResponse := NewMessageResponse{}

	bodyRequest := UpdateMessageBody{}
	reqBodyBytes := new(bytes.Buffer)

	err := json.NewEncoder(reqBodyBytes).Encode(message.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	err = json.Unmarshal(reqBodyBytes.Bytes(), &bodyRequest)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}

	newMessChan := make(chan message_service.Dto)
	g, ctx := errgroup.WithContext(e.ctx)
	g.Go(func() (err error) {
		payload := message_service.PayLoad{
			SubjectSender: message.SenderId,
			Content:       bodyRequest.Content,
			IdGroup:       bodyRequest.GroupId,
		}
		newMess, err := e.messService.UpdateMessageService(ctx, payload, bodyRequest.MessageId)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}
		_, err = e.groupService.UpdateGroupWhenHaveAction(ctx, bodyRequest.GroupId)
		if err != nil {
			fmt.Println(err)
			sentry.CaptureException(err)
			return
		}
		newMessChan <- newMess
		return
	})
	g.Go(func() (err error) {
		userOn, err := e.userOnService.GetListUSerOnlineByGroupService(ctx, bodyRequest.GroupId)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}
		select {
		case n := <-newMessChan:
			for client := range e.b.Clients {
				for _, u := range userOn {
					if u.UserID == client.UserId && u.SocketID == client.ClientId {
						if n.ParentId == int(n.ID) {
							n.ParentId = defalutParnetID
						}
						//userInfo := userdetail.GetUserById(n.SubjectSender)
						mess := Message{
							Id:                int(n.ID),
							ParentId:          n.ParentId,
							GroupId:           n.IdGroup,
							Sender:            n.SubjectSender.ID,
							Message:           n.Content,
							MessageType:       n.Type,
							TotalChildMessage: n.NumChildMess,
							CreatedAt:         n.CreatedAt,
							UpdatedAt:         n.UpdatedAt,
							DeletedAt:         n.DeletedAt,
							UserInfo:          n.SubjectSender,
						}
						bodyResponse.Message = mess
						bodyResponse.GroupId = n.IdGroup
						response.Body = bodyResponse
						msg, err = json.Marshal(response)
						if err != nil {
							return err
						}
						select {
						case client.Send <- msg:
						default:
							close(client.Send)
							delete(e.b.Clients, client)
						}
					}
				}

			}
		}
		return
	})
	if err := g.Wait(); err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}
	close(newMessChan)
}
func (e *EventImpl) DeleteMessage(message MessageRequest) {
	var msg []byte
	//var timeNow *time.Time
	//now := time.Now().UTC()

	response := MessageResponse{ResponseType: DELETED_MESSAGE}
	bodyResponse := DeleteMessageResponseBody{}

	bodyRequest := DeleteMessageBody{}
	reqBodyBytes := new(bytes.Buffer)

	err := json.NewEncoder(reqBodyBytes).Encode(message.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	err = json.Unmarshal(reqBodyBytes.Bytes(), &bodyRequest)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}

	id := bodyRequest.MessageId
	messageDelete, err := e.messService.GetMessageByIdService(e.ctx, id)

	fmt.Println(messageDelete)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	g, ctx := errgroup.WithContext(e.ctx)
	g.Go(func() (err error) {
		userOn, err := e.userOnService.GetListUSerOnlineByGroupService(ctx, bodyRequest.GroupId)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
		}
		for client := range e.b.Clients {
			for _, u := range userOn {
				if u.UserID == client.UserId && u.SocketID == client.ClientId {
					bodyResponse.MessageId = bodyRequest.MessageId
					bodyResponse.GroupId = bodyRequest.GroupId
					response.Body = bodyResponse
					msg, err = json.Marshal(response)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
					}
					select {
					case client.Send <- msg:
					default:
						close(client.Send)
						delete(e.b.Clients, client)
					}
				}

			}
		}
		return
	})
	//if messageDelete.Type == FILE_MESSAGE || messageDelete.Type == IMAGE_MESSAGE {
	//	g.Go(func() (err error) {
	//		fmt.Println("xoa tin nhan")
	//		err = e.messService.DeleteMessageByIdService(ctx, id)
	//		if err != nil {
	//			fmt.Println("loi xoa tin nhan")
	//			sentry.CaptureException(err)
	//			return
	//		}
	//		return
	//	})
	//
	//	g.Go(func() (err error) {
	//		fmt.Println("check ten group")
	//		fmt.Println(bodyRequest.GroupId)
	//		fmt.Println("----------------")
	//		err = upload.RemoveFileService(string(rune(bodyRequest.GroupId)), messageDelete.Content)
	//		if err != nil {
	//			fmt.Println("check loi ko xoa dc hình")
	//			fmt.Println(err)
	//			sentry.CaptureException(err)
	//			return
	//		}
	//		return
	//	})
	//	_, err = e.groupService.UpdateGroupWhenHaveAction(ctx, bodyRequest.GroupId)
	//	if err != nil {
	//		fmt.Println(err)
	//		sentry.CaptureException(err)
	//		return
	//	}
	//} else {
	g.Go(func() (err error) {
		err = e.messService.DeleteMessageByIdService(ctx, id)
		if err != nil {
			return
		}
		return
	})

	if messageDelete.Type == FILE_MESSAGE || messageDelete.Type == IMAGE_MESSAGE {
		g.Go(func() (err error) {
			err = upload.RemoveFileService(strconv.Itoa(bodyRequest.GroupId), messageDelete.Content)
			if err != nil {
				fmt.Printf("RemoveFileService %+v", err)
				sentry.CaptureException(err)
				return
			}
			return
		})
	}

	_, err = e.groupService.UpdateGroupWhenHaveAction(ctx, bodyRequest.GroupId)
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		return
	}
	//}
	if err := g.Wait(); err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}
	//}
}

// * 6
func (e *EventImpl) LoadChildMessage(message MessageRequest) {
	var msg []byte
	response := MessageResponse{ResponseType: MESSAGE}
	bodyResponse := NewMessageResponse{}

	bodyRequest := LoadOldMessageBody{}
	reqBodyBytes := new(bytes.Buffer)

	err := json.NewEncoder(reqBodyBytes).Encode(message.Body)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}
	err = json.Unmarshal(reqBodyBytes.Bytes(), &bodyRequest)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s \n", err)
	}

	oldMess, err := e.messService.LoadChildMessageService(e.ctx, bodyRequest.GroupId, bodyRequest.ParentMessageId, bodyRequest.LastMessageId)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}

	for client := range e.b.Clients {
		if client.UserId == message.SenderId && client.ClientId == message.ClientId && client.Identify == true {
			for _, o := range oldMess {
				if o.DeletedAt != nil {
					o.Content = ""
				} else {
					if o.Type == IMAGE_MESSAGE {
						bucketName := "group-" + strconv.Itoa(bodyRequest.GroupId)
						getLinkFile, err := upload.GetFileService(bucketName, o.Content)
						if err != nil {
							sentry.CaptureException(err)
							log.Printf("Exception : %s", err)
						}
						o.Content = getLinkFile
					}
				}
				//userInfo := userdetail.GetUserById(o.SubjectSender)
				mess := Message{
					Id:                int(o.ID),
					ParentId:          o.ParentId,
					GroupId:           o.IdGroup,
					Sender:            o.SubjectSender.ID,
					Message:           o.Content,
					MessageType:       o.Type,
					TotalChildMessage: o.NumChildMess,
					CreatedAt:         o.CreatedAt,
					UpdatedAt:         o.UpdatedAt,
					UserInfo:          o.SubjectSender,
					DeletedAt:         o.DeletedAt,
				}
				bodyResponse.GroupId = o.IdGroup
				bodyResponse.Message = mess
				response.Body = bodyResponse

				msg, err = json.Marshal(response)
				if err != nil {
					sentry.CaptureException(err)
					log.Printf("Exception : %s", err)
				}
				select {
				case client.Send <- msg:
				default:
					close(client.Send)
					delete(e.b.Clients, client)
				}
			}

		}
	}
}
