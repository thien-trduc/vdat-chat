package message

import (
	"context"
	"gitlab.com/vdat/mcsvc/chat/pkg/service/userdetail/v2"
	"time"
)

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
	Type          int
}

func (m *Messages) convertToDTO() Dto {

	message := Dto{
		ID:            m.ID,
		SubjectSender: userdetail.GetUserById(m.SubjectSender),
		Content:       m.Content,
		IdGroup:       m.IdGroup,
		ParentId:      m.ParentId,
		NumChildMess:  m.Num,
		Type:          m.Type,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
	return message
}

type Repo interface {
	GetMessagesByGroup(ctx context.Context, idGroup int) ([]Messages, error)
	GetChilMessByParentId(ctx context.Context, idGroup int, parentId int, idMess int) ([]Messages, error)
	InsertMessage(ctx context.Context, message Messages) (Messages, error)
	InsertRely(ctx context.Context, message Messages) (Messages, error)
	UpdateMessageById(ctx context.Context, message Messages, idMess int) (Messages, error)
	DeleteMessageById(ctx context.Context, idMesssage int) (err error)
	DeleteMessageByGroup(ctx context.Context, idGroup int) error
	GetContinueMessageByIdAndGroup(ctx context.Context, idMessage int, idGroup int) ([]Messages, error)
	GetMessageById(ctx context.Context, idMess int) (Messages, error)
}
type Service interface {
	AddMessageService(ctx context.Context, payload PayLoad) (Dto, error)
	AddRelyService(ctx context.Context, payload PayLoad) (Dto, error)
	LoadMessageHistoryService(ctx context.Context, idGroup int) ([]Dto, error)
	LoadChildMessageService(ctx context.Context, idGroup int, parentId int, idMess int) ([]Dto, error)
	LoadContinueMessageHistoryService(ctx context.Context, idMessage int, idGroup int) ([]Dto, error)
	DeleteMessageService(ctx context.Context, idGroup int) error
	DeleteMessageByIdService(ctx context.Context, idMess int) error
	UpdateMessageService(ctx context.Context, pay PayLoad, idMess int) (dto Dto, err error)
	GetMessageByIdService(ctx context.Context, idMess int) (Dto, error)
}
