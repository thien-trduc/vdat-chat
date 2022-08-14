import {MessageResponseType} from './message-response-type.enum';
import * as _ from 'lodash';
import {DeleteMessageResponseBody, NewMessageResponseBody, SubscribeResponseBody} from './message-response-body';

export class MessageResponse<T> {
  responseType: MessageResponseType;
  body: T;

  constructor(responseType: MessageResponseType, body: T) {
    this.responseType = responseType;
    this.body = body;
  }

  public static fromJson(data: any): MessageResponse<any> {
    if (!!data) {
      if (!_.isObject(data)) {
        data = JSON.parse(data);
      }

      const responseType: MessageResponseType = _.get(data, 'responseType', 0);
      const body: any = _.get(data, 'body', {});

      switch (responseType) {
        case MessageResponseType.SUBSCRIBED:
          return new MessageResponse<SubscribeResponseBody>(responseType, SubscribeResponseBody.fromJson(body));
        case MessageResponseType.NEW_MESSAGE:
        case MessageResponseType.MESSAGE:
        case MessageResponseType.UPDATE_MESSAGE:
          return new MessageResponse<NewMessageResponseBody>(responseType, NewMessageResponseBody.fromJson(body));
        case MessageResponseType.DELETE_MESSAGE:
          return new MessageResponse<DeleteMessageResponseBody>(responseType, DeleteMessageResponseBody.fromJson(body));
      }
    }

    return null;
  }
}
