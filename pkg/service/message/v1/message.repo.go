package message

import "context"

type Repo interface {
	GetMessagesByGroup(idChatBox int) ([]Messages, error)
	GetChilMessByParentId(idChatBox int, parentId int) ([]Messages, error)
	InsertMessage(message Messages) (Messages, error)
	InsertRely(message Messages) (Messages, error)
	//GetMessagesByChatBoxAndSeenAtOrderByCreatedAtLimit10(idChatBox int) ([]model.MessageModel, error)
	UpdateMessageById(message Messages) (Messages, error)
	//DeleteMessageById(idMesssage int) error
	//GetMessagesByGroupAndUser(idGroup int, subUser string) ([]Messages, error)
	DeleteMessageById(idMesssage int) (err error)
	DeleteMessageByGroup(idGroup int, ctx context.Context) error
	GetContinueMessageByIdAndGroup(idMessage int, idGroup int) ([]Messages, error)
}
