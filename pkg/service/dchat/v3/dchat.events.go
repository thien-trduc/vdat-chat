package dchat

type Events interface {
	SubscribeGroup(message MessageRequest)   // * 1
	LoadOldMessage(message MessageRequest)   // * 2
	SendMessage(message MessageRequest)      // * 3
	UpdateMessage(message MessageRequest)    // * 4
	DeleteMessage(message MessageRequest)    // * 5
	LoadChildMessage(message MessageRequest) // * 6
	//ReplyMessage(message MessageRequest)
}
