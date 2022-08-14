package comment

import "context"

type Repo interface {
	GetCommentById(ctx context.Context, id int64) (Comment, error)
	GetCommentByArticleID(ctx context.Context, id int64) ([]Comment, error)
	GetCommentByParentID(ctx context.Context, idParent int64) ([]Comment, error)
	InsertComment(ctx context.Context, comment Comment) (int64, error)
	InsertRelyComment(ctx context.Context, comment Comment) (int64, error)
	UpdateComment(ctx context.Context, comment Comment, id int64) error
	DeleteComment(ctx context.Context, id int64) error
}

type Service interface {
	GetCommentByArticle(ctx context.Context, idArticle int64) (results []Dto, err error)
	GetCommentByParentId(ctx context.Context, parentId int64) (results []Dto, err error)
	AddComment(ctx context.Context, payload PayLoad) (Dto, error)
	AddRelyComment(ctx context.Context, payload PayLoad) (Dto, error)
	deleteComment(ctx context.Context, idCmt int64) (err error)
	UpdateComment(ctx context.Context, payload PayLoad, id int64) (Dto, error)
}
