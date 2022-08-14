package reaction

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}

func (r *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []Reaction, err error) {
	rows, err := r.Db.QueryContext(ctx, query, args...)
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
	results = make([]Reaction, 0)
	for rows.Next() {
		r := Reaction{}
		err := rows.Scan(
			&r.ID,
			&r.IdArticle,
			&r.Type,
			&r.CreatedBy,
			&r.UpdateBy,
			&r.CreatedAt,
			&r.UpdateAt,
		)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		results = append(results, r)

	}
	rows.Close()
	return results, nil
}

func (r *RepoImpl) FindReactionByArticleUser(ctx context.Context, idArticle int64, idUser string) (results []Reaction, err error) {
	query := `select * from Reaction where id_article = $1 and create_by = $2;`
	results, err = r.fetch(ctx, query, idArticle, idUser)
	if err != nil {
		return nil, err
	}
	return
}

func (r *RepoImpl) FindById(ctx context.Context, id int64) (result Reaction, err error) {
	query := `select * from Reaction where id_reaction = $1`
	list, err := r.fetch(ctx, query, id)
	if err != nil {
		return Reaction{}, err
	}
	if len(list) > 0 {
		result = list[0]
	} else {
		return Reaction{}, ErrNotFound
	}
	return
}

func (r *RepoImpl) GetReactionByArticle(ctx context.Context, id int64) (results []Reaction, err error) {
	query := `select * from Reaction where id_article = $1;`
	results, err = r.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return
}

func (r *RepoImpl) Store(ctx context.Context, reaction *Reaction) (lastId int64, err error) {
	query := `insert into Reaction(id_article, type, create_by, update_by) values ($1,$2,$3,$4) returning id_reaction;`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx, reaction.IdArticle, reaction.Type, reaction.CreatedBy, reaction.UpdateBy).Scan(&lastId)
	if err != nil {
		return
	}
	stmt.Close()
	return
}

func (r *RepoImpl) Update(ctx context.Context, reaction *Reaction) (err error) {
	query := `update Reaction set type = $1 where id_reaction = $2;`
	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, reaction.Type, reaction.ID)
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
	return
}

func (r *RepoImpl) Delete(ctx context.Context, id int64) (err error) {
	tx, err := r.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM Reaction WHERE id_reaction = $1`, id)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return
	}
	if err = tx.Commit(); err != nil {
		log.Fatal(err)
		return
	}
	return
}
