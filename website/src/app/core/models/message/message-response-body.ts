import {Message} from './message.model';
import * as _ from 'lodash';

abstract class MessageResponseBody {
  groupId: number;

  protected constructor(groupId: number) {
    this.groupId = groupId;
  }
}

export class SubscribeResponseBody extends MessageResponseBody {
  subscribed: boolean;

  constructor(groupId: number, subscribed: boolean) {
    super(groupId);
    this.subscribed = subscribed;
  }

  public static  fromJson(body: any): SubscribeResponseBody {
    if (!!body) {
      return new SubscribeResponseBody(
        _.get(body, 'groupId', -1),
        _.get(body, 'subscribed', false)
      );
    }
  }
}

export class NewMessageResponseBody extends MessageResponseBody {
  message: Message;

  constructor(groupId: number, message: Message) {
    super(groupId);
    this.message = message;
  }

  public static fromJson(body: any): NewMessageResponseBody {
    if (!!body) {
      return new NewMessageResponseBody(
        _.get(body, 'groupId', -1),
        _.get(body, 'message', '')
      );
    }
  }
}

export class DeleteMessageResponseBody extends MessageResponseBody {
  messageId: number;

  constructor(groupId: number, messageId: number) {
    super(groupId);
    this.messageId = messageId;
  }

  public static fromJson(body: any): DeleteMessageResponseBody {
    if (!!body) {
      return new DeleteMessageResponseBody(
        _.get(body, 'groupId', -1),
        _.get(body, 'messageId', -1)
      );
    }
  }
}
