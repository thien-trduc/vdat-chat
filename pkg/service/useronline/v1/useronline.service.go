package useronline

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
)

func GetListUSerOnlineByGroupService(idGroup int) ([]Dto, error) {
	users := make([]Dto, 0)
	userOnlines, err := NewRepoImpl(database.DB).GetListUSerOnlineByGroup(idGroup)
	if err != nil {
		sentry.CaptureException(err)
		return users, err
	}
	for _, userOnline := range userOnlines {
		user := userOnline.ConvertToDto()
		users = append(users, user)
	}
	return users, nil
}
func AddUserOnlineService(payload Payload) error {
	check, err := NewRepoImpl(database.DB).GetUserOnlineBySocketIdAndHostId(payload.SocketID, payload.HostName)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	if check == (UserOnline{}) {
		online := payload.convertToModel()

		err = NewRepoImpl(database.DB).AddUserOnline(online)
		if err != nil {
			sentry.CaptureException(err)
			return err
		}
	}
	return nil
}

func DeleteUserOnlineService(socketid string, hostname string) error {
	return NewRepoImpl(database.DB).DeleteUserOnline(socketid, hostname)
}
