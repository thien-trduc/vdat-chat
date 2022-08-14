package dchat

import (
	"context"
	"encoding/json"
	"github.com/getsentry/sentry-go"
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
	ctx           context.Context
}

func NewEventImpl(b *Broker, ctx context.Context, userOnService useronline.Service, messService message_service.Service) Events {
	return &EventImpl{
		b:             b,
		userOnService: userOnService,
		messService:   messService,
		ctx:           ctx,
	}
}
func (e *EventImpl) SendMessage(message Message) {
	g, ctx := errgroup.WithContext(e.ctx)
	newMessChan := make(chan message_service.Dto)
	g.Go(func() (err error) {
		payload := message_service.PayLoad{
			SubjectSender: message.Client,
			Content:       message.Data.Body,
			IdGroup:       message.Data.GroupId,
			Type:          message.Data.Type,
		}
		newMess, err := e.messService.AddMessageService(ctx, payload)
		if err != nil {
			return
		}
		newMessChan <- newMess
		return

	})
	g.Go(func() (err error) {
		userOn, err := e.userOnService.GetListUSerOnlineByGroupService(ctx, message.Data.GroupId)
		if err != nil {
			return
		}
		select {
		case n := <-newMessChan:
			for client := range e.b.Clients {
				for _, u := range userOn {
					if u.UserID == client.UserId && u.SocketID == client.SocketId {
						if message.Data.Type == FILE || message.Data.Type == IMAGE {
							bucketName := "group-" + strconv.Itoa(message.Data.GroupId)
							getLinkFile, _ := upload.GetFileService(bucketName, message.Data.Body)
							message.Data.Body = getLinkFile
						}
						mess := Message{
							TypeEvent: SEND,
							Data: Data{
								Id:           int(n.ID),
								GroupId:      n.IdGroup,
								Body:         message.Data.Body,
								Sender:       n.SubjectSender,
								ParentID:     defalutParnetID,
								NumChildMess: n.NumChildMess,
								CreatedAt:    n.CreatedAt,
								UpdatedAt:    n.UpdatedAt,
								Type:         n.Type,
							},
						}
						msg, err := json.Marshal(mess)
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
func (e *EventImpl) SubscribeGroup(message Message) {
	historys, err := e.messService.LoadMessageHistoryService(e.ctx, message.Data.GroupId)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}

	for _, h := range historys {
		for client := range e.b.Clients {
			if client.UserId == message.Client && client.SocketId == message.Data.SocketID {
				if h.Type == FILE || h.Type == IMAGE {
					bucketName := "group-" + strconv.Itoa(message.Data.GroupId)
					getLinkFile, err := upload.GetFileService(bucketName, message.Data.Body)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
					}
					h.Content = getLinkFile
				}

				mess := Message{
					TypeEvent: SUBCRIBE,
					Data: Data{
						Id:           int(h.ID),
						GroupId:      message.Data.GroupId,
						Body:         h.Content,
						Sender:       h.SubjectSender,
						ParentID:     defalutParnetID,
						NumChildMess: h.NumChildMess,
						CreatedAt:    h.CreatedAt,
						UpdatedAt:    h.UpdatedAt,
						Type:         h.Type,
					},
				}

				msg, err := json.Marshal(mess)
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
func (e *EventImpl) ReplyMessage(message Message) {
	g, ctx := errgroup.WithContext(e.ctx)
	newMessChan := make(chan message_service.Dto)
	g.Go(func() (err error) {
		payload := message_service.PayLoad{
			SubjectSender: message.Client,
			Content:       message.Data.Body,
			IdGroup:       message.Data.GroupId,
			ID:            message.Data.Id,
			Type:          message.Data.Type,
		}
		newMess, err := e.messService.AddRelyService(ctx, payload)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
		}
		newMessChan <- newMess
		return
	})
	g.Go(func() (err error) {
		userOn, err := e.userOnService.GetListUSerOnlineByGroupService(ctx, message.Data.GroupId)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
		}
		select {
		case n := <-newMessChan:
			for client := range e.b.Clients {
				for _, u := range userOn {
					if u.UserID == client.UserId && u.SocketID == client.SocketId {
						if message.Data.Type == FILE || message.Data.Type == IMAGE {
							bucketName := "group-" + strconv.Itoa(message.Data.GroupId)
							getLinkFile, _ := upload.GetFileService(bucketName, message.Data.Body)
							message.Data.Body = getLinkFile
						}
						mess := Message{
							TypeEvent: RELY,
							Data: Data{
								Id:           int(n.ID),
								GroupId:      n.IdGroup,
								Body:         message.Data.Body,
								Sender:       n.SubjectSender,
								ParentID:     n.ParentId,
								NumChildMess: n.NumChildMess,
								CreatedAt:    n.CreatedAt,
								UpdatedAt:    n.UpdatedAt,
								Type:         n.Type,
							},
						}
						msg, err := json.Marshal(mess)
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
func (e *EventImpl) DeleteMessage(message Message) {
	id := message.Data.Id
	g, ctx := errgroup.WithContext(e.ctx)

	g.Go(func() (err error) {
		userOn, err := e.userOnService.GetListUSerOnlineByGroupService(ctx, message.Data.GroupId)
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Exception : %s", err)
		}

		for client := range e.b.Clients {
			for _, u := range userOn {
				if u.UserID == client.UserId && u.SocketID == client.SocketId {
					mess := Message{
						TypeEvent: DELETE,
						Data: Data{
							Id: id,
						},
					}
					msg, err := json.Marshal(mess)
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
		return
	})
	switch message.Data.Type {
	case TEXT:
		g.Go(func() (err error) {
			err = e.messService.DeleteMessageByIdService(ctx, id)
			if err != nil {
				return
			}
			return
		})
		break
	case FILE:
		g.Go(func() (err error) {
			err = upload.RemoveFileService(string(message.Data.GroupId), message.Data.Body)
			if err != nil {
				return
			}
			return
		})
		g.Go(func() (err error) {
			err = e.messService.DeleteMessageByIdService(ctx, id)
			if err != nil {
				return
			}
			return
		})
		break
	}
	if err := g.Wait(); err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}
}
func (e *EventImpl) LoadChildMessage(message Message) {
	child, err := e.messService.LoadChildMessageService(e.ctx, message.Data.GroupId, message.Data.Id)

	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}

	for _, h := range child {
		for client := range e.b.Clients {
			if client.UserId == message.Client && client.SocketId == message.Data.SocketID {
				if h.Type == FILE || h.Type == IMAGE {
					bucketName := "group-" + strconv.Itoa(message.Data.GroupId)
					getLinkFile, err := upload.GetFileService(bucketName, message.Data.Body)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
					}
					h.Content = getLinkFile
				}
				mess := Message{
					TypeEvent: LOADCHILDMESS,
					Data: Data{
						Id:           int(h.ID),
						GroupId:      message.Data.GroupId,
						Body:         h.Content,
						Sender:       h.SubjectSender,
						ParentID:     h.ParentId,
						NumChildMess: h.NumChildMess,
						CreatedAt:    h.CreatedAt,
						UpdatedAt:    h.UpdatedAt,
						Type:         h.Type,
					},
				}

				msg, err := json.Marshal(mess)
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
func (e *EventImpl) LoadOldMessage(message Message) {
	oldMess, err := e.messService.LoadContinueMessageHistoryService(e.ctx, message.Data.IdContinueOldMess, message.Data.GroupId)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Exception : %s", err)
	}

	for _, o := range oldMess {
		for client := range e.b.Clients {
			if client.UserId == message.Client && client.SocketId == message.Data.SocketID {
				if o.Type == FILE || o.Type == IMAGE {
					bucketName := "group-" + strconv.Itoa(message.Data.GroupId)
					getLinkFile, err := upload.GetFileService(bucketName, message.Data.Body)
					if err != nil {
						sentry.CaptureException(err)
						log.Printf("Exception : %s", err)
					}
					o.Content = getLinkFile
				}
				mess := Message{
					TypeEvent: LOADOLDMESS,
					Data: Data{
						Id:           int(o.ID),
						GroupId:      message.Data.GroupId,
						Body:         o.Content,
						Sender:       o.SubjectSender,
						ParentID:     defalutParnetID,
						NumChildMess: o.NumChildMess,
						CreatedAt:    o.CreatedAt,
						UpdatedAt:    o.UpdatedAt,
						Type:         o.Type,
					},
				}
				msg, _ := json.Marshal(mess)
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
func (e *EventImpl) UpdateMessage(message Message) {
	g, ctx := errgroup.WithContext(e.ctx)
	newMessChan := make(chan message_service.Dto)
	g.Go(func() (err error) {
		payload := message_service.PayLoad{
			SubjectSender: message.Client,
			Content:       message.Data.Body,
			IdGroup:       message.Data.GroupId,
			Type:          message.Data.Type,
		}
		newMess, err := e.messService.UpdateMessageService(ctx, payload)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}
		newMessChan <- newMess
		return
	})
	g.Go(func() (err error) {
		userOn, err := e.userOnService.GetListUSerOnlineByGroupService(ctx, message.Data.GroupId)
		if err != nil {
			sentry.CaptureException(err)
			log.Fatal(err)
		}
		select {
		case n := <-newMessChan:
			for client := range e.b.Clients {
				for _, u := range userOn {
					if u.UserID == client.UserId && u.SocketID == client.SocketId {
						if message.Data.Type == FILE || message.Data.Type == IMAGE {
							bucketName := "group-" + strconv.Itoa(message.Data.GroupId)
							getLinkFile, _ := upload.GetFileService(bucketName, message.Data.Body)
							message.Data.Body = getLinkFile
						}
						mess := Message{
							TypeEvent: SEND,
							Data: Data{
								Id:           int(n.ID),
								GroupId:      n.IdGroup,
								Body:         message.Data.Body,
								Sender:       n.SubjectSender,
								ParentID:     n.ParentId,
								NumChildMess: n.NumChildMess,
								CreatedAt:    n.CreatedAt,
								UpdatedAt:    n.UpdatedAt,
								Type:         n.Type,
							},
						}
						msg, err := json.Marshal(mess)
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
