package userdetail

import (
	"database/sql"
	"github.com/getsentry/sentry-go"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{db}
}
func (u *RepoImpl) GetListUser() ([]UserDetail, error) {
	details := make([]UserDetail, 0)
	statement := `SELECT * FROM userdetail `
	//if len(filter) > 0 {
	//	statement = statement + `WHERE fullname LIKE '` + filter + `%'`
	//}
	rows, err := u.Db.Query(statement)
	println(err)
	if err != nil {
		sentry.CaptureException(err)
		return details, err
	}
	for rows.Next() {
		var user UserDetail
		err = rows.Scan(&user.ID,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.Avatar)
		if err != nil {
			return details, err
		}
		details = append(details, user)
	}
	defer rows.Close()
	return details, nil
}
func (u *RepoImpl) AddUserDetail(detail UserDetail) error {
	statement := `insert into userdetail(user_id,role) values($1,$2)`
	_, err := u.Db.Exec(statement, detail.ID,
		detail.Role)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}
func (u *RepoImpl) GetUserDetailById(id string) (UserDetail, error) {
	var detail UserDetail
	statement := `select * from  userdetail where user_id = $1`
	rows, err := u.Db.Query(statement, id)
	if err != nil {
		sentry.CaptureException(err)
		return detail, err
	}
	if rows.Next() {
		err := rows.Scan(&detail.ID,
			&detail.Role,
			&detail.CreatedAt,
			&detail.UpdatedAt,
			&detail.DeletedAt,
			&detail.Avatar)
		if err != nil {
			return detail, err
		}
	}
	defer rows.Close()
	return detail, nil
}
func (u *RepoImpl) UpdateUserDetail(detail UserDetail) error {
	statement := `UPDATE userdetail SET role = $1  WHERE user_id = $2`
	_, err := u.Db.Exec(statement, detail.Role,
		detail.ID)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}
func (u *RepoImpl) UpdateUserById(id string, avatar string) error {
	statement := `UPDATE userdetail SET avatar = $1  WHERE user_id = $2`
	_, err := u.Db.Exec(statement, avatar, id)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}
