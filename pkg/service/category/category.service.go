package category

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
func (s *ServiceImpl) Fetch(ctx context.Context) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.Fetch(ctx)
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
func (s *ServiceImpl) GetByName(ctx context.Context, name string) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}
func (s *ServiceImpl) GetByCreatedBy(ctx context.Context, createdBy string) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.GetByCreatedBy(ctx, createdBy)
	if err != nil {
		return nil, err
	}
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}
func (s *ServiceImpl) GetByUpdateBy(ctx context.Context, updateBy string) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.GetByUpdateBy(ctx, updateBy)
	if err != nil {
		return nil, err
	}
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}
func (s *ServiceImpl) GetByParentId(ctx context.Context, parentId int64) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	list, err := s.repo.GetByParentId(ctx, parentId)
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
	m := p.convertToModel()
	m.UpdateBy = user
	m.Slug = utils.HandlerSlug(m.Name)
	err = s.repo.Update(ctx, m, id)
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
	m := p.convertToModel()
	m.CreatedBy = user
	m.UpdateBy = user
	m.Slug = utils.HandlerSlug(m.Name)
	lastId, err := s.repo.Store(ctx, m)
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
func (s *ServiceImpl) StoreChild(ctx context.Context, p *Payload, id int64, user string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	m := p.convertToModel()
	m.CreatedBy = user
	m.UpdateBy = user
	m.Slug = utils.HandlerSlug(m.Name)
	lastId, err := s.repo.StoreChild(ctx, m, id)
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
