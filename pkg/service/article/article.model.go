package article

import (
	"context"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/utils"
	"time"
)

type Article struct {
	ID         int64
	Content    string
	Title      string
	Thumbnail  string
	Version    int64
	CreatedBy  string
	UpdateBy   string
	CreatedAt  *time.Time
	UpdateAt   *time.Time
	Slug       string
	IdCategory int64
	NumReact   int64
	NumCmt     int64
	NumShare   int64
}

func (a Article) convertToDto() Dto {
	dto := Dto{
		ID:         a.ID,
		Content:    a.Content,
		Title:      a.Title,
		Thumbnail:  a.Thumbnail,
		Version:    a.Version,
		CreatedBy:  a.CreatedBy,
		UpdateBy:   a.UpdateBy,
		CreatedAt:  a.CreatedAt,
		UpdateAt:   a.UpdateAt,
		Slug:       a.Slug,
		IdCategory: a.IdCategory,
		NumReact:   a.NumReact,
		NumCmt:     a.NumCmt,
		NumShare:   a.NumShare,
	}
	return dto
}

type Repo interface {
	Fetch(ctx context.Context, pag utils.Pagination) (results []Article, err error)
	GetByID(ctx context.Context, id int64) (Article, error)
	GetByTitle(ctx context.Context, title string, pag utils.Pagination) (results []Article, err error)
	GetByUserId(ctx context.Context, userid string, pag utils.Pagination) (results []Article, err error)
	GetByUpdateBy(tx context.Context, userid string, pag utils.Pagination) (results []Article, err error)
	GetByCategoryId(ctx context.Context, idCategory int64, pag utils.Pagination) (results []Article, err error)
	Update(ctx context.Context, a *Article) error
	UpdateWithNumShare(ctx context.Context, id int64) error
	Store(ctx context.Context, a *Article) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type Service interface {
	Fetch(ctx context.Context, pag utils.Pagination) (results []Dto, err error)
	GetByID(ctx context.Context, id int64) (Dto, error)
	GetByTitle(ctx context.Context, title string, pag utils.Pagination) (results []Dto, err error)
	GetByUserId(ctx context.Context, userid string, pag utils.Pagination) (results []Dto, err error)
	GetByCategory(ctx context.Context, idCategory int64, pag utils.Pagination) (results []Dto, err error)
	Update(ctx context.Context, p *Payload, id int64, user string) (Dto, error)
	UpdateWithNumShare(ctx context.Context, id int64) (Dto, error)
	Store(ctx context.Context, p *Payload, user string) (Dto, error)
	Delete(ctx context.Context, id int64) error
}
