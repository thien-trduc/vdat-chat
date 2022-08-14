import * as _ from 'lodash';
import {IohUserAttributeLog} from './ioh-user-attribute-log';
import {IohUser} from './ioh-user.model';

export class IohUserAttribute {
  id: number;
  userId: number;
  name: string;
  value: string;
  user: IohUser;
  createdAt: Date;
  updatedAt: Date;
  aliasName?: string;
  userAttributeLogs?: Array<IohUserAttributeLog>;

  constructor(id = null, userId, name, value, user, createdAt = new Date(), updatedAt = new Date()) {
    this.id = id;
    this.userId = userId;
    this.name = name;
    this.value = value;
    this.createdAt = createdAt;
    this.updatedAt = updatedAt;
    this.user = user;
  }

  static formJson(data: any): IohUserAttribute {
    const att = new IohUserAttribute(
      data.id,
      Number(data.user_id),
      data.name,
      data.value,
      data.user,
      new Date(data.created_at),
      new Date(data.updated_at),
    );
    const userAttributeLogs = data.userAttributeLogs;
    if (userAttributeLogs && Array.isArray(userAttributeLogs)) {
      att.userAttributeLogs = userAttributeLogs.map(item => IohUserAttributeLog.fromJson(item));
    }
    return att;
  }

  static attributeToUser(attribute: IohUserAttribute): IohUser {
    const user = attribute.user;
    const att = _.cloneDeep(attribute);
    delete att.user;
    _.set(user, 'att', att);
    return user;
  }

  toJson() {
    return {
      id: this.id,
      user_id: this.userId,
      name: this.name,
      value: this.value,
      user: this.user,
      created_at: this.createdAt,
      updated_at: this.updatedAt,
    };
  }

  equal(anotherAttribute: IohUserAttribute) {
    return this.name === anotherAttribute.name && this.value === anotherAttribute.value;
  }

  setAliasName(aliasName: string) {
    this.aliasName = aliasName;
  }
}
