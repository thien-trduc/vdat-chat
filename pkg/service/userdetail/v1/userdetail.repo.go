package userdetail

import (
	"context"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
)

type Repo interface {
	GetListUser() ([]UserDetail, error)
	AddUserDetail(detail UserDetail) error
	UpdateUserDetail(etail UserDetail) error
	GetUserDetailById(id string) (UserDetail, error)
	UpdateUserById(id string, avatar string) error
}

type Repo2 interface {
	Fetch(ctx context.Context, pag utils.Pagination) (results []UserDetail, err error)
	GetListUser(ctx context.Context) ([]UserDetail, error)
	AddUserDetail(ctx context.Context, detail UserDetail) error
	UpdateUserDetail(ctx context.Context, detail UserDetail) error
	GetUserDetailById(ctx context.Context, id string) (UserDetail, error)
	UpdateUserById(ctx context.Context, id string, avatar string) error
}
