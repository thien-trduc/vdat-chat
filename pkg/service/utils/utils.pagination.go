package utils

type Pagination struct {
	Page        int32
	MaxItemPage int32
}

func (p Pagination) HandlerLimit() (offset int32) {
	offset = (p.Page - 1) * p.MaxItemPage
	return
}
