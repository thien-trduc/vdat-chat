package comment

import "time"

type Comment struct {
	ID        int64
	Content   string
	Type      string
	IdArticle int64
	ParentID  int64
	Num       int64
	Version   int64
	CreatedBy string
	UpdateBy  string
	CreatedAt *time.Time
	UpdateAt  *time.Time
}

func (c Comment) convertToDto() Dto {
	dto := Dto{
		ID:        c.ID,
		Content:   c.Content,
		Type:      c.Type,
		IdArticle: c.IdArticle,
		ParentID:  c.ParentID,
		Num:       c.Num,
		Version:   c.Version,
		CreatedBy: c.CreatedBy,
		UpdateBy:  c.UpdateBy,
		CreatedAt: c.CreatedAt,
		UpdateAt:  c.UpdateAt,
	}
	return dto
}
