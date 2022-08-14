import {MessageRequestType} from './message-request-type.enum';

export class MessageRequest<T> {
  requestType: MessageRequestType;
  clientId: string;
  senderId: string;
  body: T;

  constructor(requestType: MessageRequestType, clientId: string, senderId: string, body: T) {
    this.requestType = requestType;
    this.clientId = clientId;
    this.senderId = senderId;
    this.body = body;
  }
}
