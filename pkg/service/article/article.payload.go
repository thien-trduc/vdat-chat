package article

type Payload struct {
	Content    string `json:"content"`
	Title      string `json:"title"`
	Thumbnail  string `json:"thumbnail"`
	IdCategory int64  `json:"idCategory"`
}

func (p Payload) convertToModel() *Article {
	a := &Article{
		Content:    p.Content,
		Title:      p.Title,
		Thumbnail:  p.Thumbnail,
		IdCategory: p.IdCategory,
	}
	return a
}
