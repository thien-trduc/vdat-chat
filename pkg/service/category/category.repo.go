package category

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
func (m *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []Category, err error) {
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
	results = make([]Category, 0)
	for rows.Next() {
		c := Category{}
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.ParentID,
			&c.Num,
			&c.Version,
			&c.CreatedBy,
			&c.UpdateBy,
			&c.CreatedAt,
			&c.UpdateAt,
			&c.Slug)
		if err != nil {
			log.Panic(err)
			return nil, err
		}
		results = append(results, c)
	}
	return results, nil
}
func (m *RepoImpl) Fetch(ctx context.Context) (results []Category, err error) {
	query := `SELECT * FROM category ORDER BY created_at LIMIT 20`
	results, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByID(ctx context.Context, id int64) (result Category, err error) {
	query := `SELECT * FROM category WHERE id_category = $1`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return Category{}, err
	}
	if len(list) > 0 {
		result = list[0]
	} else {
		return Category{}, ErrNotFound
	}
	return
}
func (m *RepoImpl) GetByName(ctx context.Context, name string) (results []Category, err error) {
	query := `SELECT * FROM category WHERE name LIKE '` + name + `%' LIMIT 20`
	results, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByCreatedBy(ctx context.Context, createdBy string) (results []Category, err error) {
	query := `SELECT * FROM category WHERE create_by = $1 ORDER BY created_at LIMIT 20`
	results, err = m.fetch(ctx, query, createdBy)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByUpdateBy(ctx context.Context, updateBy string) (results []Category, err error) {
	query := `SELECT * FROM category WHERE update_by = $1 ORDER BY created_at LIMIT 20`
	results, err = m.fetch(ctx, query, updateBy)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) GetByParentId(ctx context.Context, parentId int64) (results []Category, err error) {
	query := `SELECT * FROM category WHERE parentid = $1 ORDER BY created_at LIMIT 20`
	results, err = m.fetch(ctx, query, parentId)
	if err != nil {
		return nil, err
	}
	return
}
func (m *RepoImpl) Update(ctx context.Context, a *Category, id int64) (err error) {
	query := `UPDATE category SET name = $1,parentid=$2,update_by = $3,slug = $4 WHERE id_category = $5`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, a.Name, a.ParentID, a.UpdateBy, a.Slug, id)
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
func (m *RepoImpl) Store(ctx context.Context, a *Category) (lastId int64, err error) {
	query := `INSERT INTO category(name,parentid,create_by,update_by,slug) VALUES ($1,$2,$3,$4,$5) RETURNING id_category`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx, a.Name, a.ParentID, a.CreatedBy, a.UpdateBy, a.Slug).Scan(&lastId)
	if err != nil {
		return
	}
	stmt.Close()
	return
}
func (m *RepoImpl) StoreChild(ctx context.Context, a *Category, id int64) (lastId int64, err error) {
	tx, err := m.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}
	err = tx.QueryRowContext(ctx, `INSERT INTO category(name,parentid,create_by,update_by,slug) VALUES ($1,$2,$3,$4,$5) RETURNING id_category`, a.Name, id, a.CreatedBy, a.UpdateBy, a.Slug).Scan(&lastId)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return
	}
	_, err = tx.ExecContext(ctx, `UPDATE category SET num = num + 1  WHERE id_category = $1`, id)
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
func (m *RepoImpl) Delete(ctx context.Context, id int64) (err error) {
	tx, err := m.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM article WHERE id_category = $1`, id)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return
	}
	_, err = tx.ExecContext(ctx, `DELETE FROM category WHERE id_category = $1`, id)
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
