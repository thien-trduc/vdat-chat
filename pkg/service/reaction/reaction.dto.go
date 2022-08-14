package reaction

import "time"

type Dto struct {
	ID        int64
	IdArticle int64
	Type      int8
	CreatedBy string
	UpdateBy  string
	CreatedAt *time.Time
	UpdateAt  *time.Time
}
