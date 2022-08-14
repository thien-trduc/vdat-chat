package request

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}

func (m *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []Request, err error) {
	rows, err := m.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Panic(errRow)
		}
	}()
	results = make([]Request, 0)
	for rows.Next() {
		t := Request{}
		err := rows.Scan(
			&t.ID,
			&t.IdGroup,
			&t.IdInvite,
			&t.CreatedBy,
			&t.UpdateBy,
			&t.Status,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	rows.Close()
	return results, nil
}

func (m *RepoImpl) GetAllRequest(ctx context.Context) (request []Request, err error) {
	query := `select id_request, id_Group, id_userInvite, create_by, update_by, status, created_at, updated_at from request;`
	request, err = m.fetch(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return request, err
	}
	return
}

func (m *RepoImpl) GetOneRequest(ctx context.Context, id int) (request Request, err error) {
	query := `select id_request, id_Group, id_userInvite, create_by, update_by, status, created_at, updated_at from request where id_request = $1;`
	requests, err := m.fetch(ctx, query, id)
	if err != nil {
		return Request{}, err
	}
	if len(requests) <= 0 {
		return Request{}, nil
	}
	return requests[0], nil
}

func (m *RepoImpl) GetListRequestInGroup(ctx context.Context, id int) (request []Request, err error) {
	query := `select id_request,id_Group,id_userInvite,create_by,update_by,status,created_at,updated_at
				 from request
				 WHERE id_Group = $1
				   AND status =$2;`
	request, err = m.fetch(ctx, query, id, PENDING)
	if err != nil {
		sentry.CaptureException(err)
		return request, err
	}
	return
}

func (m *RepoImpl) CreateRequest(ctx context.Context, request Request) (lastId int64, err error) {
	query := `insert into request (id_Group, id_userInvite, create_by, update_by, status) values ($1,$2,$3,$4,$5) returning id_request;`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx, request.IdGroup, request.IdInvite, request.CreatedBy, request.CreatedBy, request.Status).Scan(&lastId)
	if err != nil {
		return
	}
	stmt.Close()
	return
}

func (m *RepoImpl) UpdateRequest(ctx context.Context, id int, typeRequest int, ower string) (request Request, err error) {
	query := `update request set status = $1, updated_at = $2, update_by = $3 where id_request = $4;`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, typeRequest, time.Now(), ower, id)
	if err != nil {
		return
	}
	rowsAfected, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}
	stmt.Close()

	request, err = m.GetOneRequest(ctx, id)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	return

}

func (m *RepoImpl) GetListRequestReject(ctx context.Context, id int) (request []Request, err error) {
	query := `select id_request,id_Group,id_userInvite,create_by,update_by,status,created_at,updated_at
				 from request
				 WHERE id_Group = $1
				   AND status =$2;`
	request, err = m.fetch(ctx, query, id, REJECT)
	if err != nil {
		sentry.CaptureException(err)
		return request, err
	}
	return
}

func (m *RepoImpl) CheckExitsRequest(ctx context.Context, id string, idGroup int) (request []Request, err error) {
	query := `select id_request,id_Group,id_userInvite,create_by,update_by,status,created_at,updated_at
				 from request
				 WHERE id_Group = $1
				   AND id_userInvite =$2;`
	request, err = m.fetch(ctx, query, idGroup, id)
	if err != nil {
		sentry.CaptureException(err)
		return request, err
	}
	return
}
