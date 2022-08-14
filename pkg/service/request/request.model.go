package request

import (
	"context"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"time"
)

type AbstractModel struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
	UpdateBy  string     `json:"updateBy"`
}

type Request struct {
	AbstractModel
	CreatedBy string `json:"createdBy"`
	IdGroup   int    `json:"idGroup"`
	IdInvite  string `json:"idInvite"`
	Status    int    `json:"status"`
}

func (g *Request) ConvertToDTO() Dto {
	dto := Dto{
		Id:        g.ID,
		IdGroup:   g.IdGroup,
		IdInvite:  userdetail.GetUserById(g.IdInvite),
		Status:    g.Status,
		CreatedAt: g.CreatedAt,
		CreateBy:  userdetail.GetUserById(g.CreatedBy),
	}

	return dto
}

type Repo interface {
	GetAllRequest(ctx context.Context) (request []Request, err error)
	GetOneRequest(ctx context.Context, id int) (Request, error)
	GetListRequestInGroup(ctx context.Context, id int) (request []Request, err error)
	CreateRequest(ctx context.Context, request Request) (lastId int64, err error)
	UpdateRequest(ctx context.Context, id int, typeRequest int, ower string) (Request, error)
	GetListRequestReject(ctx context.Context, id int) (request []Request, err error)
	CheckExitsRequest(ctx context.Context, id string, idGroup int) (request []Request, err error)
}

type Service interface {
	GetAllRequest(ctx context.Context) (dto []Dto, err error)
	GetOneRequest(ctx context.Context, id int) (Dto, error)
	GetListRequestInGroup(ctx context.Context, id int) (dto []Dto, err error)
	CreateRequest(ctx context.Context, payload Payload, ower string) (dto Dto, err error)
	UpdateRequest(ctx context.Context, id int, typeRequest int, ower string) (Dto, error)
	GetListRequestReject(ctx context.Context, id int) (dto []Dto, err error)
	CheckExitsRequest(ctx context.Context, id string, idGroup int) (request Dto, err error)
}
