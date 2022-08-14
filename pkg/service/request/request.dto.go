package request

import (
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"time"
)

type Dto struct {
	Id        uint           `json:"id"`
	IdGroup   int            `json:"idGroup"`
	IdInvite  userdetail.Dto `json:"idInvite"`
	Status    int            `json:"status"`
	CreatedAt *time.Time     `json:"createdAt"`
	CreateBy  userdetail.Dto `json:"createBy"`
}
