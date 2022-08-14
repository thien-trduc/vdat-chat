package userdetail

import (
	"context"
	"fmt"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"time"
)

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

type ModelUpload struct {
	ShareUrl string
	Location string
	NameFile string
	Type     string
	CreateAt time.Time
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

func (u *UserDetail) ConvertToDto() (dto Dto, err error) {
	avatar := ""
	fmt.Println(u.Avatar)
	if u.Avatar != "" {
		avatar, err = getFileService(bucketNameAvatar, u.Avatar)
	}
	dto = Dto{
		ID:     u.ID,
		Role:   u.Role,
		Avatar: avatar,
	}
	return
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

type Repo interface {
	Fetch(ctx context.Context, pag utils.Pagination) (results []UserDetail, err error)
	GetListUser(ctx context.Context) ([]UserDetail, error)
	AddUserDetail(ctx context.Context, detail UserDetail) error
	UpdateUserDetail(ctx context.Context, detail UserDetail) error
	GetUserDetailById(ctx context.Context, id string) (UserDetail, error)
	UpdateUserById(ctx context.Context, id string, avatar string) error
}

type Service interface {
	AddUserDetailService(ctx context.Context, payload Payload) error
	UpdateUserDetailservice(ctx context.Context, payload Payload) error
	GetUserDetailByIDService(ctx context.Context, id string) (Dto, error)
	GetListUserDetailService(ctx context.Context) ([]Dto, error)
	CheckUserDetailService(ctx context.Context, payload Payload) (Dto, error)
	//GetListFromUserIdv2(ctx context.Context, listUser []string) ([]Dto,error)
}
