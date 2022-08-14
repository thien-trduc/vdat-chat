package dchat

// * Type Request Event
const (
	SUBSCRIBE          = 1
	LOAD_OLD_MESSAGE   = 2
	SEND_MESSAGE       = 3
	UPDATE_MESSAGE     = 4
	DELETE_MESSAGE     = 5
	LOAD_CHILD_MESSAGE = 6
)

type MessageRequest struct {
	RequestType int         `json:"requestType" example:"1"`
	ClientId    string      `json:"clientId"`
	SenderId    string      `json:"SenderId"" example:"user"`
	Body        interface{} `json:"body"`
}

// * 1
type SubscribeMessageBody struct {
	AccessToken string `json:"accessToken" example:1`
}

// * 2 6
type LoadOldMessageBody struct {
	GroupId         int `json:"groupId" example:1`
	ParentMessageId int `json:"parentMessageId"`
	LastMessageId   int `json:"lastMessageId"`
}

// * 3
type SendMessageBody struct {
	GroupId         int    `json:"groupId" example:1`
	ParentMessageId int    `json:"parentMessageId"`
	MessageType     int    `json:"messageType" example:1`
	Content         string `json:"content" example:"abc"`
}

// * 4
type UpdateMessageBody struct {
	GroupId   int    `json:"groupId" example:1`
	MessageId int    `json:"messageId"`
	Content   string `json:"content" example:"abc"`
}

// * 5
type DeleteMessageBody struct {
	GroupId   int `json:"groupId" example:1`
	MessageId int `json:"messageId" example:1`
}
