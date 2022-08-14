package useronline

import "time"

type Dto struct {
	HostName string     `json:"hostName"`
	SocketID string     `json:"socketId"`
	UserID   string     `json:"id"`
	LogAt    *time.Time `json:"log_at"`
}
