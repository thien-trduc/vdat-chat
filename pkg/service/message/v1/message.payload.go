package message

type PayLoad struct {
	SubjectSender string `json:"subjectSender"`
	Content       string `json:"content"`
	IdGroup       int    `json:"idGroup"`
	ID            int    `json:"id"`
	Type          string `json:"type"`
}

func (p *PayLoad) convertToModel() Messages {
	model := Messages{
		SubjectSender: p.SubjectSender,
		Content:       p.Content,
		IdGroup:       p.IdGroup,
		ParentId:      p.ID,
		Type:          p.Type,
	}
	return model
}
