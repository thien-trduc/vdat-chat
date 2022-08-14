package message

import (
	"context"
	"time"
)

type ServiceImpl struct {
	repo           Repo
	contextTimeout time.Duration
}

func NewServiceImpl(repo Repo, time time.Duration) Service {
	return &ServiceImpl{
		repo:           repo,
		contextTimeout: time,
	}
}

func (s *ServiceImpl) AddMessageService(ctx context.Context, payload PayLoad) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model := payload.convertToModel()
	m, err := s.repo.InsertMessage(ctx, model)
	if err != nil {
		return
	}
	dto = m.convertToDTO()
	return
}

func (s *ServiceImpl) AddRelyService(ctx context.Context, payload PayLoad) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model := payload.convertToModel()
	m, err := s.repo.InsertRely(ctx, model)
	if err != nil {
		return
	}
	dto = m.convertToDTO()
	return
}

func (s *ServiceImpl) LoadMessageHistoryService(ctx context.Context, idGroup int) (dtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model, err := s.repo.GetMessagesByGroup(ctx, idGroup)
	if err != nil {
		return
	}
	for _, m := range model {
		dto := m.convertToDTO()
		dtos = append(dtos, dto)
	}
	return
}

func (s *ServiceImpl) LoadChildMessageService(ctx context.Context, idGroup int, parentId int, idMess int) (dtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model, err := s.repo.GetChilMessByParentId(ctx, idGroup, parentId, idMess)
	if err != nil {
		return
	}
	for _, m := range model {
		dto := m.convertToDTO()
		dtos = append(dtos, dto)
	}
	return
}

func (s *ServiceImpl) LoadContinueMessageHistoryService(ctx context.Context, idMessage int, idGroup int) (dtos []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model, err := s.repo.GetContinueMessageByIdAndGroup(ctx, idMessage, idGroup)
	if err != nil {
		return
	}
	for _, m := range model {
		dto := m.convertToDTO()
		dtos = append(dtos, dto)
	}
	return

}

func (s *ServiceImpl) DeleteMessageService(ctx context.Context, idGroup int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err = s.repo.DeleteMessageByGroup(ctx, idGroup)
	if err != nil {
		return
	}
	return
}

func (s *ServiceImpl) DeleteMessageByIdService(ctx context.Context, idMess int) (err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	err = s.repo.DeleteMessageById(ctx, idMess)
	if err != nil {
		return
	}
	return
}

func (s *ServiceImpl) UpdateMessageService(ctx context.Context, pay PayLoad, idMess int) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model := pay.convertToModel()
	res, err := s.repo.UpdateMessageById(ctx, model, idMess)
	if err != nil {
		return
	}
	dto = res.convertToDTO()
	return
}
func (s *ServiceImpl) GetMessageByIdService(ctx context.Context, idMess int) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model, err := s.repo.GetMessageById(ctx, idMess)
	if err != nil {
		return
	}
	dto = model.convertToDTO()
	return
}
