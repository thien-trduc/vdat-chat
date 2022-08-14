package comment

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}

func (cmt RepoImpl) GetCommentById(ctx context.Context, id int64) (Comment, error) {
	cmts := make([]Comment, 0)
	query := `SELECT * FROM Comment WHERE id_cmt = $1 `

	rows, err := cmt.Db.Query(query, id)
	if err != nil {
		return cmts[0], err
	}
	for rows.Next() {
		cmt := Comment{}
		err := rows.Scan(&cmt.ID,
			&cmt.IdArticle,
			&cmt.Content,
			&cmt.Type,
			&cmt.ParentID,
			&cmt.Num,
			&cmt.Version,
			&cmt.CreatedBy,
			&cmt.UpdateBy,
			&cmt.CreatedAt,
			&cmt.UpdateAt,
		)
		if err != nil {
			return cmts[0], err
		}
		cmts = append(cmts, cmt)
	}
	defer rows.Close()
	return cmts[0], nil
}

func (cmt RepoImpl) GetCommentByArticleID(ctx context.Context, id int64) ([]Comment, error) {
	cmts := make([]Comment, 0)
	query := `SELECT * FROM Comment WHERE id_article = $1 and parentId = -1`
	rows, err := cmt.Db.Query(query, id)
	if err != nil {
		return cmts, err
	}
	for rows.Next() {
		cmt := Comment{}
		err := rows.Scan(&cmt.ID,
			&cmt.IdArticle,
			&cmt.Content,
			&cmt.Type,
			&cmt.ParentID,
			&cmt.Num,
			&cmt.Version,
			&cmt.CreatedBy,
			&cmt.UpdateBy,
			&cmt.CreatedAt,
			&cmt.UpdateAt,
		)
		if err != nil {
			return cmts, err
		}
		cmts = append(cmts, cmt)
	}
	defer rows.Close()
	return cmts, nil
}

func (cmt RepoImpl) GetCommentByParentID(ctx context.Context, idParent int64) ([]Comment, error) {
	cmts := make([]Comment, 0)
	query := `SELECT * FROM Comment WHERE parentId = $1 `
	rows, err := cmt.Db.Query(query, idParent)
	if err != nil {
		return cmts, err
	}
	for rows.Next() {
		cmt := Comment{}
		err := rows.Scan(&cmt.ID,
			&cmt.IdArticle,
			&cmt.Content,
			&cmt.Type,
			&cmt.ParentID,
			&cmt.Num,
			&cmt.Version,
			&cmt.CreatedBy,
			&cmt.UpdateBy,
			&cmt.CreatedAt,
			&cmt.UpdateAt,
		)
		if err != nil {
			return cmts, err
		}
		cmts = append(cmts, cmt)
	}
	defer rows.Close()
	return cmts, nil
}

func (cmt *RepoImpl) InsertComment(ctx context.Context, comment Comment) (lastId int64, err error) {
	statement := `INSERT INTO Comment (id_article,content,type,create_by,update_by) VALUES ($1,$2,$3,$4,$5) RETURNING id_cmt`
	stmt, err := cmt.Db.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx,
		comment.IdArticle,
		comment.Content,
		comment.Type,
		comment.CreatedBy,
		comment.UpdateBy).Scan(&lastId)

	if err != nil {
		return
	}
	stmt.Close()
	return
}

func (cmt RepoImpl) InsertRelyComment(ctx context.Context, comment Comment) (lastId int64, err error) {
	statement := `INSERT INTO Comment (id_article,content,parentId,type,create_by,update_by) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id_cmt`
	stmt, err := cmt.Db.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx,
		comment.IdArticle,
		comment.Content,
		comment.ParentID,
		comment.Type,
		comment.CreatedBy,
		comment.UpdateBy).Scan(&lastId)
	stmt.Close()

	statement1 := `update comment set num = num+1 where id_cmt = $1 `
	stmt1, err := cmt.Db.PrepareContext(ctx, statement1)
	result, _ := stmt1.ExecContext(ctx,
		comment.ParentID)
	rowsAfected, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	if err != nil {
		return
	}
	stmt1.Close()
	return
}

func (cmt RepoImpl) UpdateComment(ctx context.Context, comment Comment, id int64) (err error) {
	statement := `Update Comment set content= $1, type = $2 , version = version+1, updated_at=$3  where id_cmt = $4`
	stmt, err := cmt.Db.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, comment.Content, comment.Type, time.Now(), id)
	if err != nil {
		return err
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
	return nil

}

func (cmt RepoImpl) DeleteComment(ctx context.Context, id int64) (err error) {
	statement := `DELETE FROM Comment WHERE id_cmt=$1`
	stmt, err := cmt.Db.PrepareContext(ctx, statement)
	if err != nil {
		return
	}
	stmt.ExecContext(ctx, id)
	stmt.Close()
	statement1 := `DELETE FROM Comment WHERE parentId=$1`
	stmt1, err := cmt.Db.PrepareContext(ctx, statement1)
	if err != nil {
		return
	}
	stmt1.ExecContext(ctx, id)
	stmt1.Close()
	return nil
}
