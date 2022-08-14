package article

import (
	"context"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
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
func (s *ServiceImpl) Fetch(ctx context.Context, pag utils.Pagination) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.Fetch(ctx, pag)
	if err != nil {
		return nil, err
	}
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}
func (s *ServiceImpl) GetByID(ctx context.Context, id int64) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Dto{}, err
	}
	dto = model.convertToDto()
	return
}
func (s *ServiceImpl) GetByTitle(ctx context.Context, title string, pag utils.Pagination) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.GetByTitle(ctx, title, pag)
	if err != nil {
		return nil, err
	}
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}
func (s *ServiceImpl) GetByUserId(ctx context.Context, userid string, pag utils.Pagination) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.GetByUserId(ctx, userid, pag)
	if err != nil {
		return nil, err
	}
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}
func (s *ServiceImpl) GetByCategory(ctx context.Context, idCategory int64, pag utils.Pagination) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.GetByCategoryId(ctx, idCategory, pag)
	if err != nil {
		return nil, err
	}
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}
func (s *ServiceImpl) Update(ctx context.Context, p *Payload, id int64, user string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	a := p.convertToModel()
	a.ID = id
	a.UpdateBy = user
	a.Slug = utils.HandlerSlug(a.Title)
	err = s.repo.Update(ctx, a)
	if err != nil {
		return
	}
	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return
	}
	dto = model.convertToDto()
	return
}
func (s *ServiceImpl) UpdateWithNumShare(ctx context.Context, id int64) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	err = s.repo.UpdateWithNumShare(ctx, id)
	if err != nil {
		return
	}
	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return
	}
	dto = model.convertToDto()
	return
}
func (s *ServiceImpl) Store(ctx context.Context, p *Payload, user string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	a := p.convertToModel()
	a.CreatedBy = user
	a.UpdateBy = user
	a.Slug = utils.HandlerSlug(a.Title)
	lastId, err := s.repo.Store(ctx, a)
	if err != nil {
		return
	}
	model, err := s.repo.GetByID(ctx, lastId)
	if err != nil {
		return
	}
	dto = model.convertToDto()
	return
}
func (s *ServiceImpl) Delete(ctx context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	err = s.repo.Delete(ctx, id)
	return
}
