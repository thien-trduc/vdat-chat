package message

import "time"

type AbstractModel struct {
	ID        uint
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
type Messages struct {
	AbstractModel
	SubjectSender string
	Content       string
	IdGroup       int
	ParentId      int
	Num           int
	Type          string
}

func (m *Messages) convertToDTO() Dto {
	message := Dto{
		ID:            m.ID,
		SubjectSender: m.SubjectSender,
		Content:       m.Content,
		IdGroup:       m.IdGroup,
		ParentId:      m.ParentId,
		NumChildMess:  m.Num,
		Type:          m.Type,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	return message
}
