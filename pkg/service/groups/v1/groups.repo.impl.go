package groups

import (
	"database/sql"
	"github.com/getsentry/sentry-go"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v1"
)

type RepoImpl struct {
	Db *sql.DB
}

func NewRepoImpl(db *sql.DB) Repo {
	return &RepoImpl{Db: db}
}

const (
	thumbnail = "https://minio.nguyenchicuong.dev/public/57187617_2128229763962598_2406036693489549312_o.jpg"
)

//func (groupuser *RepoImpl) GetListGroupByUser(subUser string) ([]model.Groups, error) {
//	var groups []model.Groups
//	statement := `SELECT * FROM Groups AS g
//					INNER JOIN GroupsUsers AS g_u
//					ON g.id_group = g_u.id_group
//					WHERE g_u.sub_user_join = $1
//					ORDER BY created_at DESC
//					LIMIT 20`
//	rows, err := groupuser.Db.Query(statement, subUser)
//	if err != nil {
//		return groups, err
//	}
//	for rows.Next() {
//		groups := model.Groups{}
//		err = rows.Scan(&groups.ID, &groups.UserCreate, &groups.Name,
//			&groups.Type, &groups.CreatedAt,
//			&groups.UpdatedAt, &groups.DeletedAt)
//		if err != nil {
//			return groups, err
//		}
//		groups = append(groups, groups)
//	}
//	return groups, nil
//}
func (g *RepoImpl) GetGroupByOwnerAndUserAndTypeOne(owner string, user string) ([]Groups, error) {
	groups := make([]Groups, 0)
	statement := `SELECT g.id_group, owner_id,name, type,private,thumbnail,description, created_at, updated_at, deleted_at 
					FROM groups AS g
					INNER JOIN groups_users AS g_u
						ON g.id_group = g_u.id_group
						WHERE g.type=$1
						AND g.private=$2
						AND ((owner_id = $3 AND g_u.user_id = $4) 
							OR (owner_id = $4 AND g_u.user_id = $3))`
	rows, err := g.Db.Query(statement, ONE, true, owner, user)
	if err != nil {
		sentry.CaptureException(err)
		return groups, err
	}
	if rows.Next() {
		var group Groups
		err = rows.Scan(&group.ID,
			&group.UserCreate,
			&group.Name,
			&group.Type,
			&group.Private,
			&group.Thumbnail,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	defer rows.Close()

	return groups, nil
}
func (g *RepoImpl) GetGroupByUser(user string) ([]Groups, error) {
	groups := make([]Groups, 0)
	statement := `SELECT g.id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
 					FROM groups AS g
					INNER JOIN groups_users AS g_u
					ON g.id_group = g_u.id_group
					WHERE  g_u.user_id = $1
					ORDER BY created_at DESC 
					LIMIT 20`
	rows, err := g.Db.Query(statement, user)
	if err != nil {
		sentry.CaptureException(err)
		return groups, err
	}
	for rows.Next() {
		var group Groups
		err = rows.Scan(&group.ID,
			&group.UserCreate,
			&group.Name,
			&group.Type,
			&group.Private,
			&group.Thumbnail,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	defer rows.Close()
	return groups, nil
}
func (g *RepoImpl) GetGroupByUserAndIdGroup(user string, idGroup int) ([]Groups, error) {
	groups := make([]Groups, 0)
	statement := `SELECT g.id_group, owner_id, name, type,private,thumbnail,description,created_at, updated_at, deleted_at 
 					FROM groups AS g
					INNER JOIN groups_users AS g_u
					ON g.id_group = g_u.id_group
					WHERE  g_u.user_id = $1 and g_u.id_group = $2
					ORDER BY created_at DESC 
					LIMIT 20`
	rows, err := g.Db.Query(statement, user, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		return groups, err
	}
	for rows.Next() {
		var group Groups
		err = rows.Scan(&group.ID,
			&group.UserCreate,
			&group.Name,
			&group.Type,
			&group.Private,
			&group.Thumbnail,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	defer rows.Close()
	return groups, nil
}
func (g *RepoImpl) GetGroupByPrivateAndUser(private bool, user string) ([]Groups, error) {
	groups := make([]Groups, 0)
	statement := `((SELECT gr.id_group, owner_id,name, type,private,thumbnail,description, created_at, updated_at, deleted_at 
					FROM groups AS gr
                    WHERE gr.private = $1
					ORDER BY created_at DESC
					LIMIT 20)
            		EXCEPT
					(SELECT distinct g.id_group, owner_id,name, type,private,thumbnail,description, created_at, updated_at, deleted_at 
					from groups_users as gs inner join groups as g on gs.id_group=g.id_group  
					where user_id = $2))`
	rows, err := g.Db.Query(statement, private, user)
	if err != nil {
		sentry.CaptureException(err)
		return groups, err
	}
	for rows.Next() {
		var group Groups
		err = rows.Scan(&group.ID,
			&group.UserCreate,
			&group.Name,
			&group.Type,
			&group.Private,
			&group.Thumbnail,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	defer rows.Close()
	return groups, nil
}
func (g *RepoImpl) GetGroupByType(typeGroup string, user string) ([]Groups, error) {
	groups := make([]Groups, 0)
	statement := `SELECT *
 					FROM groups
 					WHERE type = $1
 					AND owner_id != $2
					ORDER BY created_at DESC 
					LIMIT 20`
	rows, err := g.Db.Query(statement, typeGroup, user)
	if err != nil {
		sentry.CaptureException(err)
		return groups, err
	}
	for rows.Next() {
		var group Groups
		err = rows.Scan(&group.ID,
			&group.UserCreate,
			&group.Name,
			&group.Type,
			&group.Private,
			&group.Thumbnail,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt)
		if err != nil {
			sentry.CaptureException(err)
			return groups, err
		}
		groups = append(groups, group)
	}
	defer rows.Close()
	return groups, nil
}
func (g *RepoImpl) GetOwnerByGroupAndOwner(owner string, groupId int) (bool, error) {
	statement := `SELECT owner_id FROM Groups WHERE owner_id=$1 AND id_group=$2`
	rows, err := g.Db.Query(statement, owner, groupId)
	if err != nil {
		sentry.CaptureException(err)
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	defer rows.Close()
	return false, nil
}
func (g *RepoImpl) GetGroupPublicByDoctor(user string) ([]Groups, error) {
	groups := make([]Groups, 0)
	statement := `	SELECT g.* FROM groups AS g where type = $1
					EXCEPT 
					SELECT g.* FROM groups AS g
					INNER JOIN groups_users AS g_u
					ON g.id_group = g_u.id_group
					WHERE  g_u.user_id = $2
					ORDER BY created_at DESC
					`

	rows, err := g.Db.Query(statement, MANY, user)
	if err != nil {
		return groups, err
	}
	for rows.Next() {
		var group Groups
		err = rows.Scan(&group.ID,
			&group.UserCreate,
			&group.Name,
			&group.Type,
			&group.Private,
			&group.Thumbnail,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	defer rows.Close()
	return groups, nil
}
func (g *RepoImpl) AddGroupType(group Groups) (Groups, error) {

	statement := `INSERT INTO groups (owner_id,name ,type,private,thumbnail,description) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := g.Db.Exec(statement, group.UserCreate,
		group.Name, group.Type,
		group.Private,
		thumbnail,
		group.Description)
	if err != nil {
		sentry.CaptureException(err)
		return group, err
	}

	statement = `SELECT *
 					FROM Groups AS g WHERE owner_id = $1
 					ORDER BY created_at DESC
 					LIMIT 1`
	rows, err := g.Db.Query(statement, group.UserCreate)
	if err != nil {
		sentry.CaptureException(err)
		return group, err
	}
	if rows.Next() {
		err = rows.Scan(&group.ID,
			&group.UserCreate,
			&group.Name,
			&group.Type,
			&group.Private,
			&group.Thumbnail,
			&group.Description,
			&group.CreatedAt,
			&group.UpdatedAt,
			&group.DeletedAt)
		if err != nil {
			sentry.CaptureException(err)
			return group, err
		}
	}
	defer rows.Close()
	return group, nil
}
func (g *RepoImpl) UpdateGroup(group Groups) (Groups, error) {
	var newgroup Groups
	statement := `UPDATE groups SET name=$1,description=$2,type=$3,thumbnail=$4,private=$5 WHERE id_group=$6`
	_, err := g.Db.Exec(statement, group.Name, group.Description, group.Type, group.Thumbnail, group.Private, group.ID)
	if err != nil {
		sentry.CaptureException(err)
		return newgroup, err
	}
	statement = `SELECT *
 					FROM groups
 					WHERE  id_group= $1`
	rows, err := g.Db.Query(statement, group.ID)
	if rows.Next() {
		err = rows.Scan(&newgroup.ID,
			&newgroup.UserCreate,
			&newgroup.Name,
			&newgroup.Type,
			&newgroup.Private,
			&newgroup.Thumbnail,
			&newgroup.Description,
			&newgroup.CreatedAt,
			&newgroup.UpdatedAt,
			&newgroup.DeletedAt)
		if err != nil {
			return newgroup, err
		}
	}
	defer rows.Close()
	return newgroup, nil
}
func (g *RepoImpl) DeleteGroup(idGourp int) error {
	statement := `DELETE FROM Groups_Users WHERE id_group = $1 `
	_, err := g.Db.Exec(statement, idGourp)
	if err != nil {
		return err
	}
	statement = `DELETE FROM Groups WHERE id_group = $1`
	_, err = g.Db.Exec(statement, idGourp)
	if err != nil {
		return err
	}
	return nil
}
func (g *RepoImpl) AddGroupUser(users []string, idgroup int) error {

	for _, user := range users {
		statement := `SELECT user_id FROM Groups_Users WHERE id_group=$1 AND user_id =$2`
		rows, err := g.Db.Query(statement, idgroup, user)
		if err != nil {
			sentry.CaptureException(err)
			return err
		}
		if !rows.Next() {
			statement = `INSERT INTO Groups_Users (id_group, user_id)  VALUES ($1,$2)`
			_, err := g.Db.Exec(statement, idgroup, user)
			if err != nil {
				sentry.CaptureException(err)
				return err
			}
			rows.Close()
		}

	}
	return nil
}
func (g *RepoImpl) DeleteGroupUser(users []string, idgroup int) error {
	statement := `DELETE FROM Groups_Users WHERE id_group=$1 AND user_id = $2`
	for _, user := range users {
		_, err := g.Db.Exec(statement, idgroup, user)
		if err != nil {
			sentry.CaptureException(err)
			return err
		}
	}
	return nil
}
func (g *RepoImpl) GetListUserByGroup(idGourp int) ([]userdetail.UserDetail, error) {
	users := make([]userdetail.UserDetail, 0)
	statement := `SELECT o.user_id,o.role
					FROM Groups_Users as g 
					INNER JOIN userdetail as o 
					ON g.user_id = o.user_id 					 
					WHERE id_group =$1`
	rows, err := g.Db.Query(statement, idGourp)
	if err != nil {
		sentry.CaptureException(err)
		return users, err
	}
	for rows.Next() {
		var user userdetail.UserDetail
		err = rows.Scan(&user.ID,
			&user.Role)
		if err != nil {
			sentry.CaptureException(err)
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users, nil
}
func (g *RepoImpl) GetListUserOnlineAndOfflineByGroup(idGroup int) (map[string][]userdetail.UserDetail, error) {
	mapUsers := make(map[string][]userdetail.UserDetail, 0)

	userOnlines := make([]userdetail.UserDetail, 0)
	userOffline := make([]userdetail.UserDetail, 0)

	statement := `with userInGroup as (select distinct u.user_id from groups_users as gu inner join online as u on gu.user_id = u.user_id where gu.id_group = $1)
select  u.user_id,u.role from userdetail as u inner join userInGroup as ug on u.user_id = ug.user_id`
	rows, err := g.Db.Query(statement, idGroup)

	if err != nil {
		return mapUsers, err
	}
	for rows.Next() {
		var user userdetail.UserDetail
		err = rows.Scan(&user.ID,
			&user.Role)
		if err != nil {
			sentry.CaptureException(err)
			return mapUsers, err
		}
		userOnlines = append(userOnlines, user)
	}
	rows.Close()
	statement = `with uOn as (with userInGroup as (select distinct u.user_id from groups_users as gu inner join online as u on gu.user_id = u.user_id where gu.id_group = $1)
             	select  u.user_id,u.role from userdetail as u inner join userInGroup as ug on u.user_id = ug.user_id)
				select u2.user_id,u2.role from groups_users inner join userdetail u2 on groups_users.user_id = u2.user_id where groups_users.id_group = $1
					except
				select user_id,role from uOn`

	rows, err = g.Db.Query(statement, idGroup)

	if err != nil {
		return mapUsers, err
	}
	for rows.Next() {
		var user userdetail.UserDetail
		err = rows.Scan(&user.ID,
			&user.Role)
		if err != nil {
			sentry.CaptureException(err)
			return mapUsers, err
		}
		userOffline = append(userOffline, user)
	}

	mapUsers[USERON] = userOnlines
	mapUsers[USEROFF] = userOffline
	defer rows.Close()
	return mapUsers, nil
}
