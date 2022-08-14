package useronline

import (
	"context"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"time"
)

type UserOnline struct {
	ID       string
	HostName string     `json:"hostName"`
	SocketID string     `json:"socketId"`
	UserID   string     `json:"id"`
	LogAt    *time.Time `json:"log_at"`
}

func (u *UserOnline) ConvertToDto() Dto {
	user := Dto{
		HostName: u.HostName,
		SocketID: u.SocketID,
		UserID:   u.UserID,
		LogAt:    u.LogAt,
	}
	return user
}

type Repo interface {
	Fetch(ctx context.Context, pag utils.Pagination) (results []UserOnline, err error)
	GetListUSerOnlineByGroup(ctx context.Context, idGroup int) ([]UserOnline, error)
	AddUserOnline(ctx context.Context, online UserOnline) (UserOnline, error)
	DeleteUserOnline(ctx context.Context, socketid string, hostname string) error
	GetUserOnlineBySocketIdAndHostId(ctx context.Context, socketID string, hostname string) (UserOnline, error)
}

type Service interface {
	GetListUSerOnlineByGroupService(ctx context.Context, idGroup int) ([]Dto, error)
	AddUserOnlineService(ctx context.Context, payload Payload) (Dto, error)
	DeleteUserOnlineService(ctx context.Context, socketid string, hostname string) error
}
