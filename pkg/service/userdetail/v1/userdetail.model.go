package userdetail

import "time"

type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	NotBeforePolicy  string `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
}

type UserDetail struct {
	ID        string     `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Role      string     `json:"role"`
	Avatar    string     `json:"avatar"`
}

var (
	ADMIN            = "admin"
	DOCTOR           = "doctor"
	PATIENT          = "patient"
	location         = "us-east-1"
	endpoint         = "minio.vdatlab.com"
	accessKeyID      = "gfEOZ2vrBYwoUumYEJhbcmoBLbRlonkQ"
	secretAccessKey  = "E5cGw3exgtmJVo9Q8cZgTMWJ4XNDKgjH"
	bucketNameAvatar = "avatar"
)

func (u *UserDetail) ConvertToDto() Dto {
	dto := Dto{
		ID:     u.ID,
		Role:   u.Role,
		Avatar: u.Avatar,
	}
	return dto
}

func (u *User) ConvertUserToDto() Dto {
	dto := Dto{
		ID:       u.ID,
		FullName: u.LastName + " " + u.FirstName,
		Username: u.Username,
		First:    u.FirstName,
		Last:     u.LastName,
		Role:     u.Role,
	}
	return dto
}
