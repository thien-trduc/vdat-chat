package category

type Payload struct {
	Name     string `json:"name"`
	ParentID int64  `json:"parentId"`
}

func (p Payload) convertToModel() *Category {
	category := &Category{
		Name:     p.Name,
		ParentID: p.ParentID,
	}
	return category
}
