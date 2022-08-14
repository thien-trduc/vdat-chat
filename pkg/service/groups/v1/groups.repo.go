package groups

import (
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v1"
)

type Repo interface {
	//chat one - one
	GetGroupByOwnerAndUserAndTypeOne(owner string, user string) ([]Groups, error)
	GetGroupByUser(user string) ([]Groups, error)
	GetGroupByPrivateAndUser(private bool, user string) ([]Groups, error)
	GetGroupByType(typeGroup string, user string) ([]Groups, error)
	GetOwnerByGroupAndOwner(owner string, groupId int) (bool, error)
	GetListUserByGroup(idGourp int) ([]userdetail.UserDetail, error)
	GetListUserOnlineAndOfflineByGroup(idGroup int) (map[string][]userdetail.UserDetail, error)
	GetGroupPublicByDoctor(user string) ([]Groups, error)
	AddGroupType(group Groups) (Groups, error)
	AddGroupUser(users []string, idgroup int) error
	UpdateGroup(group Groups) (Groups, error)
	DeleteGroup(idGourp int) error
	DeleteGroupUser(users []string, idgroup int) error
	GetGroupByUserAndIdGroup(user string, idGroup int) ([]Groups, error)
}
