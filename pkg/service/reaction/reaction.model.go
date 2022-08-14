package reaction

import (
	"context"
	"time"
)

type Reaction struct {
	ID        int64
	IdArticle int64
	Type      int8
	CreatedBy string
	UpdateBy  string
	CreatedAt *time.Time
	UpdateAt  *time.Time
}

func (reaction Reaction) convertToDto() Dto {
	dto := Dto{
		ID:        reaction.ID,
		Type:      reaction.Type,
		IdArticle: reaction.IdArticle,
		CreatedBy: reaction.CreatedBy,
		UpdateBy:  reaction.UpdateBy,
		CreatedAt: reaction.CreatedAt,
		UpdateAt:  reaction.UpdateAt,
	}

	return dto
}

type Repo interface {
	GetReactionByArticle(ctx context.Context, id int64) (results []Reaction, err error)
	Store(ctx context.Context, reaction *Reaction) (int64, error)
	Update(ctx context.Context, reaction *Reaction) error
	Delete(ctx context.Context, id int64) error
	FindById(ctx context.Context, id int64) (Reaction, error)
	FindReactionByArticleUser(ctx context.Context, idArticle int64, idUser string) (results []Reaction, err error)
}

type Service interface {
	GetReactionByArticle(ctx context.Context, id int64) (results []Dto, err error)
	Update(ctx context.Context, p *Payload, id int64, user string) (Dto, error)
	Store(ctx context.Context, p *Payload, user string) (Dto, error)
	Delete(ctx context.Context, id int64) error
}
