import {User} from '../user.model';
import {GroupType} from '../../constants/group-type.const';
import * as _ from 'lodash';
import {Message} from '../message/message.model';
import {CachingService} from '../../services/commons/caching.service';
import {GenerateColorService} from '../../services/commons/generate-color.service';
import {UploadFileDto} from '../file/upload-file.dto';

export class Group {
  id: number;
  nameGroup: string;
  type: GroupType;
  isPrivate: boolean;
  thumbnail: string;
  description: string;
  owner: string;
  isOwer: boolean;
  lastMessage: Message;
  isMember: boolean;
  createdAt: Date;
  updatedAt: Date;

  historyMessages: Array<Message>;
  members: Array<User> = [];
  files: Array<UploadFileDto> = [];
  images: Array<UploadFileDto> = [];
  isGroup: boolean;
  color: string;
  totalNewMessage: number;

  constructor(id: number, nameGroup: string, type: GroupType, isMember: boolean, isOwer: boolean,
              isPrivate: boolean, thumbnail: string, description: string, owner: string, lastMessage: Message,
              color?: string, createdAt: Date = new Date(), updatedAt: Date = new Date()) {
    this.id = id;
    this.nameGroup = nameGroup;
    this.type = type;
    this.isPrivate = isPrivate;
    this.thumbnail = thumbnail;
    this.lastMessage = lastMessage;
    this.owner = owner;
    this.isOwer = isOwer;
    this.description = description;
    this.isGroup = this.type === GroupType.MANY;
    this.isMember = isMember;
    this.createdAt = createdAt;
    this.updatedAt = updatedAt;

    this.members = new Array<User>();
    this.historyMessages = new Array<Message>();
    this.color = color || '#87d068';
    this.totalNewMessage = 0;
  }

  public static fromJson(data: any, generateColorService?: GenerateColorService): Group {
    const groupId: number = _.get(data, 'id', -1);
    let color = '#87d068';

    if (groupId === -1) {
      return null;
    }

    if (!!generateColorService) {
      color = generateColorService.generate(groupId.toString());
    }

    return new Group(
      groupId,
      _.get(data, 'nameGroup', ''),
      _.get(data, 'type', ''),
      _.get(data, 'isMember', false),
      _.get(data, 'isOwer', false),
      _.get(data, 'private', true),
      _.get(data, 'thumbnail', ''),
      _.get(data, 'description', ''),
      _.get(data, 'owner', ''),
      Message.fromJson(_.get(data, 'lastMessage', {})),
      color,
      new Date(_.get(data, 'createdAt', new Date().toString())),
      new Date(_.get(data, 'updatedAt', new Date().toString()))
    );
  }

  /**
   * @param message Message
   * @return is scroll to bottom
   */
  public addMessage(message: Message): boolean {
    if (this.historyMessages.length === 0) {
      this.historyMessages.push(message);
      return true;
    } else {
      const lastMessage: Message = this.historyMessages[
      this.historyMessages.length - 1
        ];
      if (lastMessage.id < message.id) {
        this.historyMessages.push(message);
        return true;
      } else {
        this.historyMessages.push(message);
        this.historyMessages = _.sortBy(this.historyMessages, 'id');
        this.historyMessages = _.uniqBy(this.historyMessages, 'id');
      }
    }

    return false;
  }

  public addMember(members: Array<User>): void {
    if (!!members && members.length > 0) {
      this.members.push(...members);
      this.members = _.uniqBy(this.members, 'userId');
    }
  }

  public caching(cachingService: CachingService): void {
    if (!!cachingService) {
      cachingService.save(this.id.toString(), this).subscribe(() => {
      });
    }
  }
}
