package article

import "time"

type Dto struct {
	ID         int64      `json:"id"`
	Content    string     `json:"content"`
	Title      string     `json:"title"`
	Thumbnail  string     `json:"thumbnail"`
	Version    int64      `json:"version"`
	CreatedBy  string     `json:"createdBy"`
	UpdateBy   string     `json:"updateBy"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdateAt   *time.Time `json:"updateAt"`
	Slug       string     `json:"slug"`
	IdCategory int64      `json:"idCategory"`
	NumReact   int64      `json:"numReact"`
	NumCmt     int64      `json:"numCmt"`
	NumShare   int64      `json:"numShare"`
}
