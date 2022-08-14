package userdetail

import (
	"context"
	"database/sql"
	"github.com/getsentry/sentry-go"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}

func (m *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []UserDetail, err error) {
	rows, err := m.Db.QueryContext(ctx, query, args...)
	if err != nil {
		sentry.CaptureException(err)
		log.Panic(err)
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			sentry.CaptureException(errRow)
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
			sentry.CaptureException(err)
			return nil, err
		}
		results = append(results, t)
	}
	rows.Close()
	return results, nil
}

func (r RepoImpl) Fetch(ctx context.Context, pag utils.Pagination) (results []UserDetail, err error) {
	query := `SELECT * FROM userdetail`
	results, err = r.fetch(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return
}

func (r RepoImpl) GetListUser(ctx context.Context) (results []UserDetail, err error) {
	query := `SELECT * FROM userdetail`
	results, err = r.fetch(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}
	return
}

func (r RepoImpl) AddUserDetail(ctx context.Context, detail UserDetail) (err error) {
	query := `insert into userdetail(user_id,role) values($1,$2)`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	_, err = stmt.ExecContext(ctx, detail.ID, detail.Role)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	stmt.Close()
	return
}

func (r RepoImpl) UpdateUserDetail(ctx context.Context, detail UserDetail) (err error) {
	query := `UPDATE userdetail SET role = $1  WHERE user_id = $2`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	_, err = stmt.ExecContext(ctx, detail.Role, detail.ID)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	stmt.Close()
	return
}

func (r RepoImpl) GetUserDetailById(ctx context.Context, id string) (result UserDetail, err error) {
	query := `select * from  userdetail where user_id = $1`
	list, err := r.fetch(ctx, query, id)
	if err != nil {
		sentry.CaptureException(err)
		return UserDetail{}, err
	}
	if len(list) > 0 {
		result = list[0]
	} else {
		return UserDetail{}, err
	}
	return
}

func (r RepoImpl) UpdateUserById(ctx context.Context, id string, avatar string) (err error) {
	query := `UPDATE userdetail SET avatar = $1  WHERE user_id = $2`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	_, err = stmt.Exec(ctx, avatar, id)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	stmt.Close()
	return
}
