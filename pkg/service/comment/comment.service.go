package comment

import (
	"context"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/database"
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

func (s *ServiceImpl) GetCommentByArticle(ctx context.Context, idArticle int64) (results []Dto, err error) {
	list, err := NewRepoImpl(database.DB).GetCommentByArticleID(ctx, idArticle)
	if err != nil {
		return nil, err
	}
	for _, cmt := range list {
		cmt.ParentID = -1
		results = append(results, cmt.convertToDto())
	}
	return results, nil
}

func (s *ServiceImpl) GetCommentByParentId(ctx context.Context, parentId int64) (results []Dto, err error) {
	list, err := NewRepoImpl(database.DB).GetCommentByParentID(ctx, parentId)
	if err != nil {
		return nil, err
	}
	for _, cmt := range list {
		results = append(results, cmt.convertToDto())
	}
	return results, nil
}

func (s *ServiceImpl) AddComment(ctx context.Context, payload PayLoad) (Dto, error) {
	var comment Comment
	cmt := payload.convertToModel()
	lastId, err := NewRepoImpl(database.DB).InsertComment(ctx, cmt)
	if err != nil {
		return comment.convertToDto(), err
	}

	newCmt, err := NewRepoImpl(database.DB).GetCommentById(ctx, lastId)

	return newCmt.convertToDto(), nil
}

func (s *ServiceImpl) AddRelyComment(ctx context.Context, payload PayLoad) (Dto, error) {
	var comment Comment
	cmt := payload.convertToModel()
	lastId, err := NewRepoImpl(database.DB).InsertRelyComment(ctx, cmt)
	if err != nil {
		return comment.convertToDto(), err
	}

	newCmt, err := NewRepoImpl(database.DB).GetCommentById(ctx, lastId)

	return newCmt.convertToDto(), nil
}

func (s *ServiceImpl) deleteComment(ctx context.Context, idCmt int64) (err error) {
	err = NewRepoImpl(database.DB).DeleteComment(ctx, idCmt)
	return err
}

func (s *ServiceImpl) UpdateComment(ctx context.Context, payload PayLoad, id int64) (Dto, error) {
	var comment Comment
	cmt := payload.convertToModel()
	err := NewRepoImpl(database.DB).UpdateComment(ctx, cmt, id)
	newCmt, err := NewRepoImpl(database.DB).GetCommentById(ctx, id)
	if err != nil {
		return comment.convertToDto(), err
	}
	return newCmt.convertToDto(), nil
}
