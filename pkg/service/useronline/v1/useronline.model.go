package useronline

import (
	"time"
)

type UserOnline struct {
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
