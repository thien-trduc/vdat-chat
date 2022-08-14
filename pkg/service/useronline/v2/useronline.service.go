package useronline

import (
	"context"
	"github.com/getsentry/sentry-go"
	"time"
)

type ServiceImpl struct {
	repo           Repo
	contextTimeout time.Duration
}

func NewServiceImpl(r Repo, time time.Duration) Service {
	return &ServiceImpl{
		repo:           r,
		contextTimeout: time,
	}
}

func (s ServiceImpl) GetListUSerOnlineByGroupService(ctx context.Context, idGroup int) ([]Dto, error) {
	users := make([]Dto, 0)
	userOnlines, err := s.repo.GetListUSerOnlineByGroup(ctx, idGroup)
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

func (s ServiceImpl) AddUserOnlineService(ctx context.Context, payload Payload) (dto Dto, err error) {
	//check, err := s.repo.GetUserOnlineBySocketIdAndHostId(ctx, payload.SocketID, payload.HostName)
	//if err != nil {
	//	sentry.CaptureException(err)
	//	return err
	//}
	online := payload.convertToModel()
	model, err := s.repo.AddUserOnline(ctx, online)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	dto = model.ConvertToDto()
	return
}

func (s ServiceImpl) DeleteUserOnlineService(ctx context.Context, socketid string, hostname string) (err error) {
	return s.repo.DeleteUserOnline(ctx, socketid, hostname)
}
