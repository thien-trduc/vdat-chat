import {MessageType} from './message-type.enum';

abstract class MessageRequestBody {
  groupId: number;
  parentMessageId: number;

  protected constructor(groupId: number, parentMessageId: number) {
    this.groupId = groupId;
    this.parentMessageId = parentMessageId;
  }
}

export class SubscribeMessageBody extends MessageRequestBody {
  accessToken: string;

  constructor(groupId: number, parentMessageId: number, accessToken: string) {
    super(groupId, parentMessageId);
    this.accessToken = accessToken;
  }
}

export class LoadOldMessageBody extends MessageRequestBody {
  lastMessageId: number;

  constructor(groupId: number, parentMessageId: number, lastMessageId: number) {
    super(groupId, parentMessageId);
    this.lastMessageId = lastMessageId;
  }
}

export class SendMessageBody extends MessageRequestBody {
  messageType: MessageType;
  content: string;

  constructor(groupId: number, parentMessageId: number, messageType: MessageType, content: string) {
    super(groupId, parentMessageId);
    this.messageType = messageType;
    this.content = content;
  }
}

export class DeleteMessageBody extends MessageRequestBody {
  messageId: number;

  constructor(groupId: number, parentMessageId: number, messageId: number) {
    super(groupId, parentMessageId);
    this.messageId = messageId;
  }
}

export class UpdateMessageBody extends MessageRequestBody {
  messageId: number;
  content: string;

  constructor(groupId: number, parentMessageId: number, messageId: number, content: string) {
    super(groupId, parentMessageId);
    this.messageId = messageId;
    this.content = content;
  }
}
