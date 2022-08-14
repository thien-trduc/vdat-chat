import * as _ from 'lodash';

export class IohUserAttributeLog {
  id: number;
  userAttributeId: number;
  userId: number;
  previousValue: any;
  newValue: any;
  createdAt: Date;
  updatedAt: Date;
  previousLabel?: string;
  newLabel?: string;

  constructor(id,
              userAttributeId,
              userId,
              previousValue,
              newValue,
              createdAt = new Date(),
              updatedAt = new Date()) {
    this.id = id;
    this.userAttributeId = userAttributeId;
    this.userId = userId;
    this.previousValue = previousValue;
    this.newValue = newValue;
    this.createdAt = createdAt;
    this.updatedAt = updatedAt;
    this.previousLabel = '';
    this.newLabel = '';
  }

  static fromJson(data): IohUserAttributeLog {
    return new IohUserAttributeLog(
      data.id,
      data.user_attribute_id,
      data.user_id,
      data.previous_value,
      data.new_value,
      new Date(data.created_at),
      new Date(data.updated_at)
    );
  }

  toJson(): IohUserAttributeLog {
    const logData = _.cloneDeep(this);

    delete logData.newLabel;
    delete logData.previousLabel;

    _.set(logData, 'user_attribute_id', logData.userAttributeId);
    _.unset(logData, 'userAttributeId');

    _.set(logData, 'user_id', logData.userId);
    _.unset(logData, 'userId');

    _.set(logData, 'previous_value', logData.previousValue);
    _.unset(logData, 'previousValue');

    _.set(logData, 'new_value', logData.newValue);
    _.unset(logData, 'newValue');

    _.set(logData, 'created_at', logData.createdAt);
    _.unset(logData, 'createdAt');

    _.set(logData, 'updated_at', logData.updatedAt);
    _.unset(logData, 'updatedAt');

    return logData;
  }

  setLabel(previousLabel = '', newLabel = '') {
    this.previousLabel = previousLabel;
    this.newLabel = newLabel;
    return this;
  }

  getPreviousValue() {
    return _.get(this, 'previousValue.value', '');
  }

  getNewValue() {
    return _.get(this, 'newValue.value', '');
  }
}
