package comment

type PayLoad struct {
	Content   string `json:"content"`
	Type      string `json:"type"`
	UserId    string `json:"userId"`
	ParentID  int64  `json:"parentId"`
	IdArticle int64  `json:"idArticle"`
}

func (p PayLoad) convertToModel() Comment {
	c := Comment{
		Content:   p.Content,
		Type:      p.Type,
		IdArticle: p.IdArticle,
		ParentID:  p.ParentID,
		CreatedBy: p.UserId,
		UpdateBy:  p.UserId,
	}
	return c
}
