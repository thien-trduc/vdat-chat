package category

import (
	"context"
	"time"
)

type Category struct {
	ID        int64
	Name      string
	ParentID  int64
	Num       int64
	Version   int64
	CreatedBy string
	UpdateBy  string
	CreatedAt *time.Time
	UpdateAt  *time.Time
	Slug      string
}

func (category Category) convertToDto() Dto {
	dto := Dto{
		ID:        category.ID,
		Name:      category.Name,
		ParentID:  category.ParentID,
		Num:       category.Num,
		Version:   category.Version,
		CreatedBy: category.CreatedBy,
		UpdateBy:  category.UpdateBy,
		CreatedAt: category.CreatedAt,
		UpdateAt:  category.UpdateAt,
		Slug:      category.Slug,
	}
	return dto
}

type Repo interface {
	Fetch(ctx context.Context) (results []Category, err error)
	GetByID(ctx context.Context, id int64) (Category, error)
	GetByName(ctx context.Context, name string) (results []Category, err error)
	GetByCreatedBy(ctx context.Context, createdBy string) (results []Category, err error)
	GetByUpdateBy(ctx context.Context, updateBy string) (results []Category, err error)
	GetByParentId(ctx context.Context, parentId int64) (results []Category, err error)
	Update(ctx context.Context, a *Category, id int64) error
	Store(ctx context.Context, a *Category) (int64, error)
	StoreChild(ctx context.Context, a *Category, id int64) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type Service interface {
	Fetch(ctx context.Context) (results []Dto, err error)
	GetByID(ctx context.Context, id int64) (Dto, error)
	GetByName(ctx context.Context, name string) (results []Dto, err error)
	GetByCreatedBy(ctx context.Context, createdBy string) (results []Dto, err error)
	GetByUpdateBy(ctx context.Context, updateBy string) (results []Dto, err error)
	GetByParentId(ctx context.Context, parentId int64) (results []Dto, err error)
	Update(ctx context.Context, p *Payload, id int64, user string) (Dto, error)
	Store(ctx context.Context, p *Payload, user string) (Dto, error)
	StoreChild(ctx context.Context, p *Payload, id int64, user string) (Dto, error)
	Delete(ctx context.Context, id int64) error
}
