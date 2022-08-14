import * as _ from 'lodash';
import {RequestStatus} from '../../constants/request-status.enum';
import {User} from "../user.model";
import {GenerateColorService} from "../../services/commons/generate-color.service";

export class RequestDto {
  requestId: number;
  groupId: number;
  userInvite: User;
  userInvited: User;
  createdAt: Date;
  status: RequestStatus;

  constructor(requestId: number, groupId: number, userInvited: User,
              createdAt: Date, status: RequestStatus, userInvite?: User) {
    this.requestId = requestId;
    this.groupId = groupId;
    this.userInvite = userInvite;
    this.userInvited = userInvited;
    this.createdAt = createdAt;
    this.status = status;
  }

  public static fromJson(data: any, generateColorService?: GenerateColorService): RequestDto {
    return new RequestDto(
      _.get(data, 'id', -1),
      _.get(data, 'idGroup', -1),
      User.fromJson(_.get(data, 'idInvite', {}), generateColorService),
      new Date(_.get(data, 'createdAt', new Date().toString())),
      _.get(data, 'status', RequestStatus.PENDING),
      User.fromJson(_.get(data, 'createBy', {}), generateColorService),
    );
  }
}
