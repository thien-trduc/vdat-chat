package userdetail

import (
	"context"
	"database/sql"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
)

type RepoImpl2 struct {
	Db *sql.DB
}

func (m *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []UserDetail, err error) {
	rows, err := m.Db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Panic(errRow)
		}
	}()
	results = make([]UserDetail, 0)
	for rows.Next() {
		t := UserDetail{}
		err := rows.Scan(
			&t.ID,
			&t.Role,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
			&t.Avatar,
		)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		results = append(results, t)
	}
	rows.Close()
	return results, nil
}

func (r RepoImpl2) Fetch(ctx context.Context, pag utils.Pagination) (results []UserDetail, err error) {
	panic("implement me")
}

func (r RepoImpl2) GetListUser(ctx context.Context) ([]UserDetail, error) {
	panic("implement me")
}

func (r RepoImpl2) AddUserDetail(ctx context.Context, detail UserDetail) error {
	panic("implement me")
}

func (r RepoImpl2) UpdateUserDetail(ctx context.Context, detail UserDetail) error {
	panic("implement me")
}

func (r RepoImpl2) GetUserDetailById(ctx context.Context, id string) (UserDetail, error) {
	panic("implement me")
}

func (r RepoImpl2) UpdateUserById(ctx context.Context, id string, avatar string) error {
	panic("implement me")
}

func NewRepoImpl2(db *sql.DB) Repo2 {
	return &RepoImpl2{db}
}
