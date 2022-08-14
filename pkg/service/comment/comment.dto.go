package comment

import "time"

type Dto struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Type      string     `json:"type"`
	IdArticle int64      `json:"idArticle"`
	ParentID  int64      `json:"parentId"`
	Num       int64      `json:"num"`
	Version   int64      `json:"version"`
	CreatedBy string     `json:"createdBy"`
	UpdateBy  string     `json:"updateBy"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdateAt  *time.Time `json:"updateAt"`
}
