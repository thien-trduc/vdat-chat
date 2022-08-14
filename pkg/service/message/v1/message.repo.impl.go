package message

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}

func (mess *RepoImpl) GetMessagesByGroup(idChatBox int) ([]Messages, error) {
	messages := make([]Messages, 0)
	statement := `SELECT id_mess,user_sender,content,id_group,numChild,created_at,updated_at,type FROM messages WHERE id_group = $1 and parentID IS NULL ORDER BY created_at DESC LIMIT 20`
	rows, err := mess.Db.Query(statement, idChatBox)
	if err != nil {
		sentry.CaptureException(err)
		return messages, err
	}
	for rows.Next() {
		m := Messages{}
		err := rows.Scan(&m.ID,
			&m.SubjectSender,
			&m.Content,
			&m.IdGroup,
			&m.Num,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Type,
		)
		if err != nil {
			sentry.CaptureException(err)
			return messages, err
		}
		messages = append(messages, m)
	}
	defer rows.Close()
	return messages, nil
}

func (mess *RepoImpl) GetChilMessByParentId(idChatBox int, parentId int) ([]Messages, error) {
	messages := make([]Messages, 0)
	statement := `SELECT id_mess,user_sender,content,id_group,parentID,numChild,created_at,updated_at,type FROM messages WHERE id_group = $1 and parentID = $2 ORDER BY created_at DESC LIMIT 20`
	rows, err := mess.Db.Query(statement, idChatBox, parentId)
	if err != nil {
		sentry.CaptureException(err)
		return messages, err
	}
	for rows.Next() {
		m := Messages{}
		err := rows.Scan(&m.ID,
			&m.SubjectSender,
			&m.Content,
			&m.IdGroup,
			&m.ParentId,
			&m.Num,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Type,
		)
		if err != nil {
			sentry.CaptureException(err)
			return messages, err
		}
		messages = append(messages, m)
	}
	defer rows.Close()
	return messages, nil
}

func (mess *RepoImpl) InsertMessage(message Messages) (Messages, error) {
	fmt.Println(message.SubjectSender)
	fmt.Println(message.Content)
	fmt.Println(message.IdGroup)
	var id int
	m := Messages{}
	statement := `INSERT INTO messages (user_sender,content,id_group,type) VALUES ($1,$2,$3,$4) RETURNING id_mess`
	err := mess.Db.QueryRow(statement,
		message.SubjectSender,
		message.Content,
		message.IdGroup,
		message.Type).Scan(&id)

	statement = `SELECT id_mess,user_sender,content,id_group,numChild,created_at,updated_at FROM messages WHERE  id_mess = $1`
	rows, err := mess.Db.Query(statement, id)
	if rows.Next() {
		err = rows.Scan(&m.ID,
			&m.SubjectSender,
			&m.Content,
			&m.IdGroup,
			&m.Num,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			sentry.CaptureException(err)
			return m, err
		}
	}
	defer rows.Close()
	return m, err
}
func (mess *RepoImpl) InsertRely(message Messages) (Messages, error) {
	var id int
	m := Messages{}
	statement := `INSERT INTO messages (user_sender,content,id_group,parentID,type) VALUES ($1,$2,$3,$4,$5) RETURNING id_mess`
	err := mess.Db.QueryRow(statement,
		message.SubjectSender,
		message.Content,
		message.IdGroup,
		message.ParentId,
		message.Type).Scan(&id)
	statement = `SELECT id_mess,user_sender,content,id_group,parentID,numChild,created_at,updated_at FROM messages WHERE  id_mess = $1`
	rows, err := mess.Db.Query(statement, id)
	if rows.Next() {
		err = rows.Scan(&m.ID,
			&m.SubjectSender,
			&m.Content,
			&m.IdGroup,
			&m.ParentId,
			&m.Num,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			sentry.CaptureException(err)
			return m, err
		}
		defer rows.Close()
	}
	statement1 := `update messages set numchild = numchild+1 where id_mess = $1 RETURNING id_mess`
	_ = mess.Db.QueryRow(statement1,
		message.ParentId).Scan(&id)
	return m, err
}

//func (mess *MessageRepoImpl) GetMessagesByChatBoxAndSeenAtOrderByCreatedAtLimit10(idChatBox int) ([]model.MessageModel, error) {
//	message := make([]model.MessageModel, 0)
//	statement := `SELECT * FROM message WHERE id_chat = $1 AND seen_at IS NULL ORDER BY created_at LIMIT 10`
//	rows, err := mess.Db.Query(statement, idChatBox)
//	if err != nil {
//		return message, err
//	}
//	for rows.Next() {
//		message := model.MessageModel{}
//		err := rows.Scan(&message.ID, &message.IdChat, &message.Content, &message.SeenAt, &message.CreatedAt, &message.UpdatedAt, &message.DeletedAt)
//		if err != nil {
//			return message, err
//		}
//		message = append(message, message)
//	}
//	return message, nil
//}
func (mess *RepoImpl) UpdateMessageById(message Messages) (result Messages, err error) {
	statement := `UPDATE messages SET content = $1 where id_mess = $2 and id_group = $3`
	_, err = mess.Db.Exec(statement, message.Content, message.ID, message.IdGroup)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	statement = `SELECT id_mess,user_sender,content,id_group,numChild,created_at,updated_at,type FROM messages WHERE id_mess=$1`
	rows, err := mess.Db.Query(statement, message.ID)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	if rows.Next() {
		err = rows.Scan(&result.ID,
			&result.SubjectSender,
			&result.Content,
			&result.IdGroup,
			&result.Num,
			&result.CreatedAt,
			&result.UpdatedAt,
			&result.Type,
		)
	}
	defer rows.Close()
	return
}

//func (mess *RepoImpl) GetMessagesByGroupAndUser(idGroup int, subUser string) ([]Messages, error) {
//	messages := make([]Messages, 0)
//	statement := `SELECT m.id_mess,m.subject_sender,m.content,m.created_at,m.updated_at,m.deleted_at,m.id_group
//					FROM (SELECT * FROM dchat.public.message WHERE id_group = $1) AS m
//					LEFT JOIN Messages_Delete AS md
//					ON m.id_mess = md.id_mess
//					WHERE  m.created_at >= (SELECT g_u.last_deleted_messages
//											FROM Groups AS g
//											INNER JOIN Groups_Users AS g_u
//											ON g.id_group = $1
//											AND g.id_group = g_u.id_group
//											WHERE g_u.sub_user_join = $2
//											ORDER BY g_u.last_deleted_messages DESC
//											LIMIT 1)
//					AND (md.sub_user_deleted != $2 OR md.sub_user_deleted IS NULL)
//					`
//	rows, err := mess.Db.Query(statement, idGroup, subUser)
//	if err != nil {
//		return messages, err
//	}
//	for rows.Next() {
//		message := Messages{}
//		err := rows.Scan(&message.ID,
//			&message.SubjectSender,
//			&message.Content,
//			&message.CreatedAt,
//			&message.UpdatedAt,
//			&message.DeletedAt,
//			&message.IdGroup)
//		if err != nil {
//			return messages, err
//		}
//		messages = append(messages, message)
//	}
//	defer rows.Close()
//	return messages, nil
//}

//func (mess *MessageRepoImpl) DeleteMessageById(idMesssage int) error {
//	statement := `DELETE FROM message WHERE id_mess = $1`
//	_, err := mess.Db.Exec(statement, idMesssage)
//	return err
//}
func (mess *RepoImpl) DeleteMessageByGroup(idGroup int, ctx context.Context) error {
	s := `DELETE FROM messages WHERE id_group = $1`
	tx, err := mess.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}
	_, execErr := tx.Exec(s, idGroup)
	if execErr != nil {
		_ = tx.Rollback()
		log.Fatal(execErr)
		return execErr
	}
	if err := tx.Commit(); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
		return err
	}
	return nil
}
func (mess *RepoImpl) GetContinueMessageByIdAndGroup(idMessage int, idGroup int) ([]Messages, error) {
	messages := make([]Messages, 0)
	statement := `select id_mess,user_sender,content,id_group,numChild,created_at,updated_at,type from messages where id_mess < $1 and id_group = $2 and parentID IS NULL order by created_at DESC limit 20`
	rows, err := mess.Db.Query(statement, idMessage, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		return messages, err
	}
	for rows.Next() {
		m := Messages{}
		err := rows.Scan(&m.ID,
			&m.SubjectSender,
			&m.Content,
			&m.IdGroup,
			&m.Num,
			&m.CreatedAt,
			&m.UpdatedAt,
			&m.Type,
		)
		if err != nil {
			sentry.CaptureException(err)
			return messages, err
		}
		messages = append(messages, m)
	}
	defer rows.Close()
	return messages, nil
}
func (mess *RepoImpl) DeleteMessageById(idMesssage int) (err error) {
	s := `DELETE FROM messages WHERE id_mess = $1`
	_, err = mess.Db.Exec(s, idMesssage)
	if err != nil {
		sentry.CaptureException(err)
		log.Print(err)
		return
	}
	return
}
