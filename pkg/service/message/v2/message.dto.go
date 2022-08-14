package message

import (
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"time"
)

type Dto struct {
	ID            uint           `json:"id"`
	SubjectSender userdetail.Dto `json:"subjectSender"`
	Content       string         `json:"message"`
	IdGroup       int            `json:"idGroup"`
	ParentId      int            `json:"parentId"`
	NumChildMess  int            `json:"numChildMess"`
	Type          int            `json:"messageType"`
	CreatedAt     *time.Time     `json:"createdAt"`
	UpdatedAt     *time.Time     `json:"updatedAt"`
	DeletedAt     *time.Time     `json:"deletedAt"`
}
