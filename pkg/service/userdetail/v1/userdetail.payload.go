package userdetail

type Payload struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
	Username string `json:"userName"`
	First    string `json:"first"`
	Last     string `json:"last"`
	Role     string `json:"role"`
}

func (p *Payload) convertToModel() UserDetail {
	u := UserDetail{
		ID:   p.ID,
		Role: p.Role,
	}
	return u
}
