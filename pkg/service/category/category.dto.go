package category

import "time"

type Dto struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	ParentID  int64      `json:"parentId"`
	Num       int64      `json:"num"`
	Version   int64      `json:"version"`
	CreatedBy string     `json:"createdBy"`
	UpdateBy  string     `json:"updateBy"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdateAt  *time.Time `json:"updateAt"`
	Slug      string     `json:"slug"`
}
