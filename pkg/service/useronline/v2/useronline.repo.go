package useronline

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

func (r *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []UserOnline, err error) {
	rows, err := r.Db.QueryContext(ctx, query, args...)
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
	results = make([]UserOnline, 0)
	for rows.Next() {
		t := UserOnline{}
		err := rows.Scan(
			&t.ID,
			&t.HostName,
			&t.SocketID,
			&t.UserID,
			&t.LogAt,
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

func (r RepoImpl) Fetch(ctx context.Context, pag utils.Pagination) (results []UserOnline, err error) {
	panic("implement me")
}

func (r RepoImpl) GetListUSerOnlineByGroup(ctx context.Context, idGroup int) ([]UserOnline, error) {
	var users []UserOnline
	statement := `select u.user_id,u.socket_id from groups_users as gu inner join online as u on gu.user_id = u.user_id where gu.id_group = $1`
	rows, err := r.Db.Query(statement, idGroup)
	//println(err)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user UserOnline
		err = rows.Scan(&user.UserID, &user.SocketID)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users, nil
}

func (r RepoImpl) AddUserOnline(ctx context.Context, online UserOnline) (o UserOnline, err error) {
	var lastId int64
	query := `INSERT INTO ONLINE (hostname,socket_id,user_id) VALUES ($1,uuid_in(md5(random()::text || clock_timestamp()::text)::cstring),$2) RETURNING id_onl`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	err = stmt.QueryRowContext(ctx, online.HostName, online.UserID).Scan(&lastId)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	stmt.Close()
	query = `select * from online where id_onl = $1`
	results, err := r.fetch(ctx, query, lastId)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	if len(results) > 0 {
		o = results[0]
		return
	}
	return
}

func (r RepoImpl) DeleteUserOnline(ctx context.Context, socketid string, hostname string) (err error) {
	query := `DELETE FROM ONLINE WHERE socket_id=$1 AND hostname=$2`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	_, err = stmt.ExecContext(ctx, socketid, hostname)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	stmt.Close()
	return
}

func (r RepoImpl) GetUserOnlineBySocketIdAndHostId(ctx context.Context, socketID string, hostname string) (result UserOnline, err error) {
	query := `SELECT hostname,socket_id,user_id,log_at FROM ONLINE WHERE hostname=$1 AND socket_id=$2`
	results, err := r.fetch(ctx, query, hostname, socketID)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	return results[0], err
}
