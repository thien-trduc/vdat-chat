package groups

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"golang.org/x/sync/errgroup"
	"log"
	"strconv"
	"time"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}

func (m *RepoImpl) fetch(ctx context.Context, query string, args ...interface{}) (results []Groups, err error) {
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
	results = make([]Groups, 0)
	for rows.Next() {
		t := Groups{}
		err := rows.Scan(
			&t.ID,
			&t.UserCreate,
			&t.Name,
			&t.Type,
			&t.Private,
			&t.Thumbnail,
			&t.Description,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	rows.Close()
	return results, nil
}

func (m *RepoImpl) GetGroupByNameForDoctor(ctx context.Context, user string, keyword string, pag utils.Pagination) (groups []Groups, err error) {
	query := `select id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
				from groups
				where (type = $1 and deleted_at IS NULL and name ILIKE $2) or id_group in (
					select g.id_group
					from groups g left join groups_users gu on g.id_group = gu.id_group
					where g.name ILIKE $2
                          and gu.id_group is not null
                          and user_id = $3
                          and g.deleted_at IS NULL)
				group by id_group
				order by updated_at DESC OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	groups, err = m.fetch(ctx, query, MANY, `%`+keyword+`%`, user)
	if err != nil {
		return
	}
	return
}

func (m *RepoImpl) GetGroupByNameForPatient(ctx context.Context, private bool, user string, keyword string, pag utils.Pagination) (groups []Groups, err error) {
	query := `select g.id_group, owner_id,name, type,private,thumbnail,description, created_at, updated_at, deleted_at
				from groups
				where (private = $1 and deleted_at IS NULL and name ILIKE $2) or id_group in (
					select g.id_group
					from groups g left join groups_users gu on g.id_group = gu.id_group
					where g.name ILIKE $2
					  and gu.id_group is not null
					  and user_id = $3
					  and g.deleted_at IS NULL
					)
				group by id_group
				order by updated_at DESC OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	groups, err = m.fetch(ctx, query, private, `%`+keyword+`%`, user)
	if err != nil {
		return
	}
	return
}

func (m *RepoImpl) GetGroupByOwnerAndUserAndTypeOne(ctx context.Context, owner string, user string) (groups []Groups, err error) {
	query := `SELECT g.id_group, owner_id,name, type,private,thumbnail,description, created_at, updated_at, deleted_at 
					FROM groups AS g
					INNER JOIN groups_users AS g_u
					ON g.id_group = g_u.id_group
					WHERE g.type=$1 
							AND deleted_at IS NULL
							AND g.private=$2 
							AND ((owner_id = $3 AND g_u.user_id = $4) 
							OR (owner_id = $4 AND g_u.user_id = $3))`
	groups, err = m.fetch(ctx, query, ONE, true, owner, user)
	if err != nil {
		return
	}
	return
}
func (m *RepoImpl) GetGroupByUser(ctx context.Context, user string) (groups []Groups, err error) {
	query := `SELECT g.id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
					FROM groups AS g
					INNER JOIN groups_users AS g_u
					ON g.id_group = g_u.id_group
					WHERE  g_u.user_id = $1 and g.deleted_at IS NULL
					ORDER BY updated_at DESC 
					LIMIT 20`
	groups, err = m.fetch(ctx, query, user)
	if err != nil {
		return
	}
	return
}
func (m *RepoImpl) GetGroupByUserAndIdGroup(ctx context.Context, user string, idGroup int) (groups []Groups, err error) {
	query := `SELECT g.id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
					FROM groups AS g
					INNER JOIN groups_users AS g_u
					ON g.id_group = g_u.id_group
					WHERE  g_u.user_id = $1 and g_u.id_group = $2 and g.deleted_at IS NULL
					ORDER BY updated_at DESC 
					LIMIT 20`
	groups, err = m.fetch(ctx, query, user, idGroup)
	if err != nil {
		return
	}
	return
}

func (m *RepoImpl) GetGroupIdGroup(ctx context.Context, idGroup int) (group Groups, err error) {
	query := `SELECT id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
					FROM groups 
					WHERE  id_group = $1`
	groups, err := m.fetch(ctx, query, idGroup)
	if err != nil {
		return
	}
	if len(groups) <= 0 {
		return group, nil
	}
	return groups[0], nil
}

func (m *RepoImpl) GetGroupByPrivateAndUser(ctx context.Context, private bool, user string, pag utils.Pagination) (groups []Groups, err error) {
	query := `select id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
			from groups
			where (private = $1 and deleted_at IS NULL) or id_group in (
				select g.id_group
				from groups g left join groups_users gu on g.id_group = gu.id_group
				where gu.id_group is not null and user_id = $2 and g.deleted_at IS NULL
				)
			group by id_group
			order by updated_at DESC OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	groups, err = m.fetch(ctx, query, private, user)
	if err != nil {
		return
	}
	return
}
func (m *RepoImpl) GetGroupByType(ctx context.Context, typeGroup string, user string) (groups []Groups, err error) {
	query := `SELECT * FROM groups WHERE type = $1 AND owner_id != $2 and deleted_at IS NULL ORDER BY updated_at DESC LIMIT 20`
	groups, err = m.fetch(ctx, query, typeGroup, user)
	if err != nil {
		return
	}
	return
}
func (m *RepoImpl) GetOwnerByGroupAndOwner(ctx context.Context, owner string, groupId int) (check bool, err error) {
	check = false
	query := `SELECT * FROM Groups WHERE owner_id=$1 AND id_group=$2 and deleted_at IS NULL`
	groups, err := m.fetch(ctx, query, owner, groupId)
	if err != nil {
		return
	}
	if len(groups) > 0 {
		check = true
		return
	}
	return
}
func (m *RepoImpl) GetGroupPublicByDoctor(ctx context.Context, user string, pag utils.Pagination) (groups []Groups, err error) {
	query := `select id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
				from groups
				where (type = $1 and deleted_at IS NULL) or id_group in (
					select g.id_group
					from groups g left join groups_users gu on g.id_group = gu.id_group
					where gu.id_group is not null and user_id = $2 and g.deleted_at IS NULL)
				group by id_group
				order by updated_at DESC OFFSET ` + strconv.Itoa(int(pag.HandlerLimit())) + ` LIMIT ` + strconv.Itoa(int(pag.MaxItemPage))
	groups, err = m.fetch(ctx, query, MANY, user)
	if err != nil {
		return
	}
	return
}
func (m *RepoImpl) AddGroupType(ctx context.Context, group Groups) (model Groups, err error) {
	var lastId int64
	query := `INSERT INTO groups (owner_id,name ,type,private,thumbnail,description) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id_group`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	err = stmt.QueryRowContext(ctx, group.UserCreate, group.Name, group.Type, group.Private, thumbnail, group.Description).Scan(&lastId)
	if err != nil {
		return
	}
	stmt.Close()

	query = `SELECT * FROM Groups AS g WHERE id_group = $1`
	groups, err := m.fetch(ctx, query, lastId)
	if err != nil {
		return
	}
	if len(groups) > 0 {
		model = groups[0]
		return
	}
	return
}
func (m *RepoImpl) UpdateGroup(ctx context.Context, group Groups) (model Groups, err error) {
	query := `UPDATE groups SET name=$1,
								description=$2,
								private=$3 
							WHERE id_group=$4`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, group.Name,
		group.Description,
		group.Private,
		group.ID)
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

	query = `SELECT * FROM groups WHERE  id_group= $1`
	groups, err := m.fetch(ctx, query, group.ID)
	if err != nil {
		return
	}
	if len(groups) > 0 {
		model = groups[0]
		return
	}
	return
}

func (m *RepoImpl) UpdateGroupWithThumbnail(ctx context.Context, group Groups) (model Groups, err error) {
	query := `UPDATE groups SET name=$1,
								description=$2,
								thumbnail=$3,
								private=$4 
							WHERE id_group=$5`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, group.Name,
		group.Description,
		group.Thumbnail,
		group.Private,
		group.ID)
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

	query = `SELECT * FROM groups WHERE  id_group= $1`
	groups, err := m.fetch(ctx, query, group.ID)
	if err != nil {
		return
	}
	if len(groups) > 0 {
		model = groups[0]
		return
	}
	return
}

func (m *RepoImpl) DeleteGroup(ctx context.Context, idGourp int) (err error) {
	tx, err := m.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM messages WHERE id_group = $1`, idGourp)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM Groups_Users WHERE id_group = $1`, idGourp)
	if err != nil {
		_ = tx.Rollback()
		return
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM Groups WHERE id_group = $1`, idGourp)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	return

}
func (m *RepoImpl) AddGroupUser(ctx context.Context, users []string, idgroup int) error {
	tx, err := m.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	for _, user := range users {
		query := `SELECT user_id FROM Groups_Users WHERE id_group=$1 AND user_id =$2`
		rows, err := m.Db.QueryContext(ctx, query, idgroup, user)
		if err != nil {
			return err
		}
		if !rows.Next() {
			query = `INSERT INTO Groups_Users (id_group, user_id)  VALUES ($1,$2)`
			_, err = tx.ExecContext(ctx, query, idgroup, user)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			rows.Close()
		}
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}
func (m *RepoImpl) DeleteGroupUser(ctx context.Context, users []string, idgroup int) error {
	tx, err := m.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	query := `DELETE FROM Groups_Users WHERE id_group=$1 AND user_id = $2`
	for _, user := range users {
		_, err = tx.ExecContext(ctx, query, idgroup, user)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}
func (m *RepoImpl) GetListUserByGroup(ctx context.Context, idGourp int) (users []userdetail.UserDetail, err error) {
	query := `SELECT o.user_id,o.role
					FROM Groups_Users as g
					INNER JOIN userdetail as o
					ON g.user_id = o.user_id
					WHERE id_group =$1`
	rows, err := m.Db.QueryContext(ctx, query, idGourp)
	if err != nil {
		return
	}
	for rows.Next() {
		var user userdetail.UserDetail
		err = rows.Scan(&user.ID,
			&user.Role)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	defer rows.Close()
	return
}
func (m *RepoImpl) GetListUserOnlineAndOfflineByGroup(ctx context.Context, idGroup int) (mapUsers map[string][]userdetail.UserDetail, err error) {
	userOnlines := make([]userdetail.UserDetail, 0)
	userOffline := make([]userdetail.UserDetail, 0)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		query := `with userInGroup as (select distinct u.user_id from groups_users as gu inner join online as u on gu.user_id = u.user_id where gu.id_group = $1)
select u.user_id,u.role from userdetail as u inner join userInGroup as ug on u.user_id = ug.user_id`
		rows, err := m.Db.QueryContext(ctx, query, idGroup)
		fmt.Println(err)
		if err != nil {
			return
		}

		for rows.Next() {
			var user userdetail.UserDetail
			err = rows.Scan(&user.ID,
				&user.Role)
			fmt.Sprintln(user)
			if err != nil {
				return
			}
			userOnlines = append(userOnlines, user)
		}
		rows.Close()
		return
	})
	g.Go(func() (err error) {
		query := `with uOn as (with userInGroup as (select distinct u.user_id from groups_users as gu inner join online as u on gu.user_id = u.user_id where gu.id_group = $1)
            	select  u.user_id,u.role from userdetail as u inner join userInGroup as ug on u.user_id = ug.user_id)
				select u2.user_id,u2.role from groups_users inner join userdetail u2 on groups_users.user_id = u2.user_id where groups_users.id_group = $1
					except
				select user_id,role from uOn`

		rows, err := m.Db.QueryContext(ctx, query, idGroup)

		if err != nil {
			panic(err)
			return
		}
		for rows.Next() {
			var user userdetail.UserDetail
			fmt.Sprintln(user)
			err = rows.Scan(&user.ID,
				&user.Role)
			if err != nil {
				return
			}
			userOffline = append(userOffline, user)
		}
		defer rows.Close()
		return
	})
	if err = g.Wait(); err != nil {
		return
	}

	mapUsers = make(map[string][]userdetail.UserDetail, 0)
	mapUsers[USERON] = userOnlines
	mapUsers[USEROFF] = userOffline
	fmt.Println(mapUsers)
	return
}

func (m *RepoImpl) TempDeleteGroup(ctx context.Context, idGroup int) (err error) {
	query := `UPDATE groups SET deleted_at = now() WHERE id_group=$1`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, idGroup)
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

func (m *RepoImpl) GetGroupNeedDelete(ctx context.Context) (groups []Groups, err error) {
	query := `Select * from Groups WHERE deleted_at IS NOT NULL;`
	groups, err = m.fetch(ctx, query)
	if err != nil {
		return
	}
	return
}

func (m *RepoImpl) UpdateGroupWhenHaveAction(ctx context.Context, idGroup int) (model Groups, err error) {
	query := `update groups set updated_at = $1 where id_group = $2`
	stmt, err := m.Db.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	result, err := stmt.ExecContext(ctx, time.Now(), idGroup)
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

	query = `SELECT * FROM groups WHERE  id_group= $1`
	groups, err := m.fetch(ctx, query, idGroup)
	if err != nil {
		return
	}
	if len(groups) > 0 {
		model = groups[0]
		return
	}
	return
}
