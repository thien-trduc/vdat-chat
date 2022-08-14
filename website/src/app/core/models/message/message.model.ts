import {User} from '../user.model';
import {MessageType} from './message-type.enum';
import {Group} from '../group/group.model';
import * as _ from 'lodash';

export class Message {
  id: number;
  parentId: number;
  groupId: number;
  sender: User;
  message: string;
  messageType: MessageType;
  totalChildMessage: number;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date;
  createdBy: string;
  updatedBy: string;
  group: Group;

  replyMessages: Array<Message>;
  isOwner: boolean;

  constructor(
    id: number,
    groupId: number,
    sender: User,
    messageType: MessageType,
    totalChildMessage: number = 0,
    message: any,
    parentId: number,
    createdAt: Date,
    updatedAt: Date,
    deletedAt: Date
  ) {
    this.id = id;
    this.messageType = messageType;
    this.totalChildMessage = totalChildMessage;
    this.message = message;
    this.parentId = parentId;
    this.groupId = groupId;
    this.createdAt = createdAt;
    this.updatedAt = updatedAt;
    this.deletedAt = deletedAt;
    this.sender = sender;

    this.group = Group.fromJson({id: groupId});
    this.replyMessages = new Array<Message>();
    this.isOwner = false;
  }

  public static fromJson(data: any): Message {
    if (!!data) {
      const createdAt = _.get(data, 'createdAt', null);
      const updatedAt = _.get(data, 'updatedAt', null);
      const deletedAt = _.get(data, 'deletedAt', null);

      return new Message(
        _.get(data, 'id', -1),
        _.get(data, 'groupId', -1),
        User.fromJson(_.get(data, 'userInfo', null)),
        _.get(data, 'messageType', MessageType.TEXT_MESSAGE),
        _.get(data, 'totalChildMessage', 0),
        _.get(data, 'message', ''),
        _.get(data, 'parentId', -1),
        !!createdAt ? new Date(createdAt) : new Date(),
        !!updatedAt ? new Date(updatedAt) : new Date(),
        !!deletedAt ? new Date(deletedAt) : null,
      );
    }

    return null;
  }
}
