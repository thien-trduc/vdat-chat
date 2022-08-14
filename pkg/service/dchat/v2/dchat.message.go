package dchat

import (
	"time"
)

type Message struct {
	TypeEvent string `json:"type" example:"send_text"`
	Data      Data   `json:"data" `
	Client    string
}
type Data struct {
	GroupId           int        `json:"groupId" example:1`
	Id                int        `json:"id" example=1`
	Body              string     `json:"body" example:"tin nhan moi"`
	Sender            string     `json:"sender" example:"abc"`
	SocketID          string     `json:"socketId" example:"9999"`
	IdContinueOldMess int        `json:"idContinueOldMess"`
	ParentID          int        `json:"parentID"`
	Type              string     `json:"type"`
	NumChildMess      int        `json:"numChildMess"`
	Status            string     `json:"status" example:"done"`
	CreatedAt         *time.Time `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt"`
}
type Events interface {
	SendMessage(message Message)
	SubscribeGroup(message Message)
	ReplyMessage(message Message)
	DeleteMessage(message Message)
	LoadChildMessage(message Message)
	LoadOldMessage(message Message)
	UpdateMessage(message Message)
}

type ResponseHistoryMess struct {
	Historys []Message `json:"historys"`
}

//func (d *Data) prepare(newMess message_service.Dto){
//
//}
