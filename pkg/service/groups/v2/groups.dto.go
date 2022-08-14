package groups

import (
	"encoding/json"
	"fmt"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/message/v2"
	"time"
)

type Dto struct {
	Id          uint        `json:"id"`
	Name        string      `json:"nameGroup"`
	Type        string      `json:"type"`
	Private     bool        `json:"private"`
	Owner       string      `json:"owner"`
	Thumbnail   string      `json:"thumbnail"`
	Description string      `json:"description"`
	IsMember    bool        `json:"isMember"`
	IsOwer      bool        `json:"isOwer"`
	LastMessage message.Dto `json:"lastMessage"`
	CreatedAt   *time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time  `json:"updatedAt"`
}
type Dtos []Dto

func (d Dto) MarshalToJsonString() string {
	b, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}
