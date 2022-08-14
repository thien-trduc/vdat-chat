package reaction

import (
	"context"
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

func (s *ServiceImpl) GetReactionByArticle(ctx context.Context, id int64) (results []Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	list, err := s.repo.GetReactionByArticle(ctx, id)
	for _, a := range list {
		dto := a.convertToDto()
		results = append(results, dto)
	}
	return
}

func (s *ServiceImpl) Update(ctx context.Context, p *Payload, id int64, user string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	reaction := p.convertToModel()
	reaction.ID = id
	reaction.UpdateBy = user
	err = s.repo.Update(ctx, reaction)
	if err != nil {
		return
	}
	model, err := s.repo.FindById(ctx, id)
	if err != nil {
		return
	}
	dto = model.convertToDto()
	return
}

func (s *ServiceImpl) Store(ctx context.Context, p *Payload, user string) (dto Dto, err error) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()
	reaction := p.convertToModel()
	reaction.CreatedBy = user
	reaction.UpdateBy = user

	checkReaction, _ := s.repo.FindReactionByArticleUser(ctx, reaction.IdArticle, user)
	if len(checkReaction) > 0 && reaction.Type != checkReaction[0].Type {
		reaction.ID = checkReaction[0].ID
		_ = s.repo.Update(ctx, reaction)
		model, _ := s.repo.FindById(ctx, checkReaction[0].ID)
		dto = model.convertToDto()
		return
	}
	if len(checkReaction) > 0 {
		err = s.repo.Delete(ctx, checkReaction[0].ID)
		return
	} else {
		lastId, _ := s.repo.Store(ctx, reaction)
		model, _ := s.repo.FindById(ctx, lastId)
		dto = model.convertToDto()
		return
	}
}

func (s *ServiceImpl) Delete(ctx context.Context, id int64) error {
	panic("implement me")
}
