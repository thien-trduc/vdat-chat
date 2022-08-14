import {EncryptService} from '../services/commons/encrypt.service';
import * as _ from 'lodash';

export class Invite {
  groupId: number;
  groupName: string;
  userInviteId: string;
  userInviteFullName: string;
  timeInvite: Date;

  constructor(groupId: number, groupName: string, userInviteId: string,
              userInviteFullName: string, timeInvite: Date) {
    this.groupId = groupId;
    this.groupName = groupName;
    this.userInviteId = userInviteId;
    this.userInviteFullName = userInviteFullName;
    this.timeInvite = timeInvite;
  }

  public static fromJson(data: any): Invite {
    return new Invite(
      _.get(data, 'groupId', -1),
      _.get(data, 'groupName', -1),
      _.get(data, 'userInviteId', ''),
      _.get(data, 'userInviteFullName', ''),
      new Date(_.get(data, 'timeInvite', new Date().toString())),
    );
  }

  public static fromData(data: string, encryptService: EncryptService): Invite {
    return JSON.parse(encryptService.decrypt(data));
  }

  public toData(encryptService: EncryptService): string {
    return encryptService.encrypt(JSON.stringify(this));
  }
}
