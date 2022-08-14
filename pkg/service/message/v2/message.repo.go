package message

import (
	"context"
	"database/sql"
	"github.com/getsentry/sentry-go"
	"log"
	"time"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{
		Db: db,
	}
}
func (r *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []Messages, err error) {
	rows, err := r.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Panic(errRow)
		}
	}()
	results = make([]Messages, 0)
	for rows.Next() {
		t := Messages{}
		err := rows.Scan(
			&t.ID,
			&t.SubjectSender,
			&t.Content,
			&t.IdGroup,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
			&t.ParentId,
			&t.Num,
			&t.Type,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	rows.Close()
	return results, nil
}

func (r *RepoImpl) GetMessagesByGroup(ctx context.Context, idGroup int) (messages []Messages, err error) {
	//query := `SELECT * FROM messages WHERE id_group = $1 and parentID IS NULL ORDER BY created_at DESC LIMIT 20`
	query := `SELECT * FROM messages WHERE id_group = $1 and parentID = id_mess ORDER BY created_at DESC LIMIT 20`
	messages, err = r.fetch(ctx, query, idGroup)
	if err != nil {
		return
	}
	return
}

func (r *RepoImpl) GetChilMessByParentId(ctx context.Context, idGroup int, parentId int, idMess int) (messages []Messages, err error) {
	query := `WITH msGroup AS (SELECT * FROM messages WHERE id_group = $1 and parentID != id_mess)
				SELECT * FROM msGroup WHERE parentID = $2 `
	check := `and id_mess < $3 `
	sort := `ORDER BY created_at DESC LIMIT 20`
	messages, err = r.fetch(ctx, query, idGroup, parentId, idMess)
	if idMess != 0 {
		query = query + check + sort
		messages, err = r.fetch(ctx, query, idGroup, parentId, idMess)
	} else {
		query = query + sort
		messages, err = r.fetch(ctx, query, idGroup, parentId)
	}
	if err != nil {
		return
	}
	return
}

func (r *RepoImpl) InsertMessage(ctx context.Context, message Messages) (m Messages, err error) {
	var idNew int
	tx, err := r.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	query := `INSERT INTO messages (user_sender,content,id_group,type) VALUES ($1,$2,$3,$4) RETURNING id_mess`
	err = tx.QueryRowContext(ctx, query, message.SubjectSender, message.Content, message.IdGroup, message.Type).Scan(&idNew)
	if err != nil {
		_ = tx.Rollback
		return
	}
	query = `UPDATE messages SET parentid = $1 WHERE id_mess=$2`
	_, err = tx.ExecContext(ctx, query, idNew, idNew)
	if err != nil {
		_ = tx.Rollback
		return
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback
		return
	}
	query = `SELECT * FROM messages WHERE id_mess = $1`
	messages, err := r.fetch(ctx, query, idNew)
	if err != nil {
		return
	}
	if len(messages) > 0 {
		m = messages[0]
	}
	return
}
func (r *RepoImpl) InsertRely(ctx context.Context, message Messages) (m Messages, err error) {
	var idNew int
	tx, err := r.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})

	query := `INSERT INTO messages (user_sender,content,id_group,parentID,type) VALUES ($1,$2,$3,$4,$5) RETURNING id_mess`
	err = tx.QueryRowContext(ctx, query, message.SubjectSender, message.Content, message.IdGroup, message.ParentId, message.Type).Scan(&idNew)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	query = `UPDATE messages SET numchild = numchild+1 WHERE id_mess = $1`
	_, err = tx.ExecContext(ctx, query, message.ParentId)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return
	}

	query = `SELECT * FROM messages WHERE  id_mess = $1`
	messages, err := r.fetch(ctx, query, idNew)
	if err != nil {
		return
	}
	if len(messages) > 0 {
		m = messages[0]
	}
	return
}

func (r *RepoImpl) UpdateMessageById(ctx context.Context, message Messages, idMess int) (m Messages, err error) {
	query := `UPDATE messages SET content = $1 where id_mess = $2 and id_group = $3`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, message.Content, idMess, message.IdGroup)
	if err != nil {
		return
	}
	stmt.Close()

	query = `SELECT * FROM messages WHERE id_mess=$1`
	messages, err := r.fetch(ctx, query, idMess)
	if err != nil {
		return
	}
	if len(messages) > 0 {
		m = messages[0]
	}
	return
}

func (r *RepoImpl) DeleteMessageByGroup(ctx context.Context, idGroup int) error {
	s := `DELETE FROM messages WHERE id_group = $1`
	tx, err := r.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sentry.CaptureException(err)

	}
	_, execErr := tx.Exec(s, idGroup)
	if execErr != nil {
		_ = tx.Rollback()

		return execErr
	}
	if err := tx.Commit(); err != nil {
		sentry.CaptureException(err)

		return err
	}
	return nil
}

func (r *RepoImpl) GetContinueMessageByIdAndGroup(ctx context.Context, idMessage int, idGroup int) (messages []Messages, err error) {
	query := `with msGroup as (select * from messages where id_group = $1 and parentID = id_mess ) 
				select * from msGroup `
	check := `where id_mess < $2 `
	sort := `order by created_at DESC limit 20`
	if idMessage != 0 {
		query = query + check + sort
		messages, err = r.fetch(ctx, query, idGroup, idMessage)
	} else {
		query = query + sort
		messages, err = r.fetch(ctx, query, idGroup)
	}
	if err != nil {
		return
	}
	return
}

func (r *RepoImpl) DeleteMessageById(ctx context.Context, idMesssage int) (err error) {
	query := `UPDATE messages SET deleted_at = $1 WHERE id_mess = $2`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	_, err = stmt.ExecContext(ctx, time.Now(), idMesssage)
	if err != nil {
		return
	}
	return
}

func (r *RepoImpl) GetMessageById(ctx context.Context, idMess int) (message Messages, err error) {
	query := `select * from messages where id_mess = $1`
	messages, err := r.fetch(ctx, query, idMess)
	if err != nil {
		return
	}
	if len(messages) > 0 {
		message = messages[0]
	}
	return
}
