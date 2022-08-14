package useronline

type Repo interface {
	GetListUSerOnlineByGroup(idGroup int) ([]UserOnline, error)
	AddUserOnline(online UserOnline) error
	DeleteUserOnline(socketid string, hostname string) error
	GetUserOnlineBySocketIdAndHostId(socketID string, hostname string) (UserOnline, error)
}
