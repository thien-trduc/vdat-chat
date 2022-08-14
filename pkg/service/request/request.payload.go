package request

type Payload struct {
	IdInvite string `json:"idInvite"`
	IdGroup  int    `json:"idGroup"`
	CreateBy string
}

func (p Payload) ConvertToModel() Request {
	r := Request{
		IdInvite:  p.IdInvite,
		IdGroup:   p.IdGroup,
		CreatedBy: p.CreateBy,
		Status:    PENDING,
	}
	return r
}
