package request

import (
	"context"
	"fmt"
	"github.com/getsentry/sentry-go"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/groups/v2"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"golang.org/x/sync/errgroup"
	"time"
)

type ServiceImpl struct {
	repo           Repo
	userService    userdetail.Service
	groupService   groups.Service
	contextTimeout time.Duration
}

func NewServiceImpl(r Repo, userService userdetail.Service, groupService groups.Service, time time.Duration) Service {
	return &ServiceImpl{
		repo:           r,
		userService:    userService,
		groupService:   groupService,
		contextTimeout: time,
	}
}

func (s *ServiceImpl) GetAllRequest(ctx context.Context) (dtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		requests, err := s.repo.GetAllRequest(ctx)
		fmt.Println(requests)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		for _, request := range requests {
			dto := request.ConvertToDTO()
			dtos = append(dtos, dto)
		}

		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	return
}

func (s *ServiceImpl) GetOneRequest(ctx context.Context, id int) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		request, err := s.repo.GetOneRequest(ctx, id)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		dto = request.ConvertToDTO()
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	return
}

func (s *ServiceImpl) GetListRequestInGroup(ctx context.Context, id int) (dtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		requests, err := s.repo.GetListRequestInGroup(ctx, id)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		for _, request := range requests {
			dto := request.ConvertToDTO()
			dtos = append(dtos, dto)
		}
		return
	})

	if err = g.Wait(); err != nil {
		return
	}
	return
}

func (s *ServiceImpl) CreateRequest(ctx context.Context, payload Payload, ower string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	fmt.Println(payload)

	checkRole, err := s.groupService.CheckRoleOwnerInGroupService(ctx, ower, payload.IdGroup)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	if checkRole {
		fmt.Println("la chu")
		User := []string{payload.IdInvite}
		_ = s.groupService.AddUserInGroupService(ctx, User, payload.IdGroup)
		dto.Status = APPROVE
		return
	}

	requestConvert := payload.ConvertToModel()
	fmt.Println("check")
	fmt.Println(requestConvert)
	index, err := s.repo.CreateRequest(ctx, requestConvert)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	request, err := s.repo.GetOneRequest(ctx, int(index))
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	dto = request.ConvertToDTO()
	return
}

func (s *ServiceImpl) UpdateRequest(ctx context.Context, id int, typeRequest int, ower string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	request, err := s.repo.UpdateRequest(ctx, id, typeRequest, ower)
	if typeRequest == APPROVE {
		User := []string{request.IdInvite}
		_ = s.groupService.AddUserInGroupService(ctx, User, request.IdGroup)
	}
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	dto = request.ConvertToDTO()
	return
}

func (s *ServiceImpl) GetListRequestReject(ctx context.Context, id int) (dtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() (err error) {
		requests, err := s.repo.GetListRequestReject(ctx, id)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
		for _, request := range requests {
			dto := request.ConvertToDTO()
			dtos = append(dtos, dto)
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	return
}

func (s *ServiceImpl) CheckExitsRequest(ctx context.Context, id string, idGroup int) (request Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	requests, err := s.repo.CheckExitsRequest(ctx, id, idGroup)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
	if len(requests) > 0 {
		request = requests[0].ConvertToDTO()
	} else {
		request = Dto{}
	}
	return
}
