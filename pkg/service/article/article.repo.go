package article

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"log"
	"strconv"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}
func (m *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []Article, err error) {
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
	results = make([]Article, 0)
	for rows.Next() {
		t := Article{}
		err := rows.Scan(
			&t.ID,
			&t.Content,
			&t.Title,
			&t.Thumbnail,
			&t.Version,
			&t.CreatedBy,
			&t.UpdateBy,
			&t.CreatedAt,
			&t.UpdateAt,
			&t.Slug,
			&t.IdCategory,
			&t.NumReact,
			&t.NumCmt,
			&t.NumShare,
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
func (m *RepoImpl) Fetch(ctx context.Context, pag utils.Pagination) (results []Article, err error) {
	query := `SELECT * FROM article ORDER BY created_at OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	results, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByID(ctx context.Context, id int64) (result Article, err error) {
	query := `SELECT * FROM article WHERE id_article = $1 `
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return Article{}, err
	}
	if len(list) > 0 {
		result = list[0]
	} else {
		return Article{}, ErrNotFound
	}
	return
}
func (m *RepoImpl) GetByTitle(ctx context.Context, title string, pag utils.Pagination) (results []Article, err error) {
	query := `SELECT * FROM article WHERE title LIKE '%` + title + `%' OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	results, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByUserId(ctx context.Context, userid string, pag utils.Pagination) (results []Article, err error) {
	query := `SELECT * FROM article WHERE create_by = $1 ORDER BY created_at OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	results, err = m.fetch(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByCategoryId(ctx context.Context, idCategory int64, pag utils.Pagination) (results []Article, err error) {
	query := `SELECT * FROM article WHERE id_category = $1 ORDER BY created_at OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	results, err = m.fetch(ctx, query, idCategory)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByUpdateBy(ctx context.Context, userid string, pag utils.Pagination) (results []Article, err error) {
	query := `SELECT * FROM article WHERE update_by = $1 ORDER BY created_at OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	results, err = m.fetch(ctx, query, userid)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) Update(ctx context.Context, a *Article) (err error) {
	query := `UPDATE article SET title = $1,content = $2,thumbnail = $3,update_by = $4,slug = $5,id_category = $6 WHERE id_article = $7`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Thumbnail, a.UpdateBy, a.Slug, a.IdCategory, a.ID)
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
func (m *RepoImpl) UpdateWithNumShare(ctx context.Context, id int64) (err error) {
	query := `UPDATE article SET num_share = num_share + 1 WHERE id_article = $1`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, id)
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
func (m *RepoImpl) Store(ctx context.Context, a *Article) (lastId int64, err error) {
	query := `INSERT INTO article(title,content,thumbnail,create_by,update_by,slug,id_category) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id_article`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx, a.Title, a.Content, a.Thumbnail, a.CreatedBy, a.UpdateBy, a.Slug, a.IdCategory).Scan(&lastId)
	if err != nil {
		return
	}
	stmt.Close()
	return
}
func (m *RepoImpl) Delete(ctx context.Context, id int64) (err error) {
	tx, err := m.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM comment WHERE id_article = $1`, id)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM article WHERE id_article = $1`, id)
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
