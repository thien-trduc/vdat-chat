package dchat

import (
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"time"
)

// * Type Respone Event
const (
	SUBSCRIBED      = 1
	MESSAGE         = 2
	NEW_MESSAGE     = 3
	UPDATED_MESSAGE = 4
	DELETED_MESSAGE = 5
)

type MessageResponse struct {
	ResponseType int         `json:"responseType" example:1`
	Body         interface{} `json:"body"`
}
type Message struct {
	Id                int            `json:"id" example=1`
	ParentId          int            `json:"parentId"`
	GroupId           int            `json:"groupId" example:1`
	Sender            string         `json:"sender" example:"abc"`
	Message           string         `json:"message" example:"tin nhan moi"`
	MessageType       int            `json:"messageType"`
	TotalChildMessage int            `json:"totalChildMessage"`
	CreatedAt         *time.Time     `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         *time.Time     `json:"deletedAt"`
	UserInfo          userdetail.Dto `json:"userInfo"`
}
type SubscribeResponseBody struct {
	Subscribed bool `json:"subscribed"`
}
type DeleteMessageResponseBody struct {
	GroupId   int `json:"groupId"`
	MessageId int `json:"messageId"`
}
type NewMessageResponse struct {
	GroupId int     `json:"groupId"`
	Message Message `json:"message"`
}
