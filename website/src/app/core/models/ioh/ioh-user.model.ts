import * as _ from 'lodash';
import * as moment from 'moment';
import {IohUserAttribute} from './ioh-user-attribute';
import {DateTimeFormatInfo} from '../../constants/date-time-format.info';

export class IohUser {
  public id: number;
  public female: boolean;
  public email: string;
  public phone: string;
  public firstName: string;
  public lastName: string;
  public address: string;
  public birthday: Date;
  public wardId: string;
  public createdAt: Date;
  public updatedAt: Date;
  public isActive: boolean;
  public metadata: string;
  public profileId: string;
  public patientId: string;
  public fullName: string;
  public uuid: string;
  public config: any;
  public userAttribute?: IohUserAttribute;

  constructor(id, email, phone, firstName, lastName,
              address, birthday, female = true, wardId, createdAt,
              updatedAt, isActive) {
    this.id = id;
    this.female = female;
    this.email = email;
    this.phone = phone;
    this.firstName = firstName;
    this.lastName = lastName;
    this.address = address;
    this.birthday = birthday;
    this.wardId = wardId;
    this.createdAt = createdAt ? createdAt : new Date();
    this.updatedAt = updatedAt ? updatedAt : new Date();
    this.isActive = isActive;
    this.metadata = ' ';
  }

  static fromJson(data) {
    const user = new IohUser(
      parseInt(data.id, null),
      data.email,
      data.phone,
      data.first_name,
      data.last_name,
      data.address,
      new Date(data.birthday),
      data.female,
      data.ward_id,
      data.created_at,
      data.updated_at,
      data.is_active);
    user.uuid = data.uuid;
    user.config = data.config;
    user.profileId = data.profile_id;
    user.patientId = data.patient_id;
    user.fullName = IohUser.getFullName(user);
    return user;
  }

  static getFullName(user: IohUser, location: string = 'vi') {
    if (location === 'vi') {
      return `${user.lastName} ${user.firstName}`;
    }
    return `${user.firstName} ${user.lastName}`;
  }

  toJson() {
    const clone = _.cloneDeep(this);

    _.set(clone, 'first_name', this.firstName);
    _.unset(clone, 'firstName');

    _.set(clone, 'last_name', this.lastName);
    _.unset(clone, 'lastName');

    _.set(clone, 'ward_id', this.wardId);
    _.unset(clone, 'wardId');

    _.set(clone, 'created_at', this.createdAt);
    _.unset(clone, 'createdAt');

    _.set(clone, 'updated_at', this.updatedAt);
    _.unset(clone, 'updatedAt');

    _.set(clone, 'is_active', this.isActive);
    _.unset(clone, 'isActive');

    _.set(clone, 'profile_id', this.profileId);
    _.unset(clone, 'profileId');

    _.set(clone, 'patient_id', this.patientId);
    _.unset(clone, 'patientId');
    _.unset(clone, 'fullName');
    _.unset(clone, 'userAttribute');

    _.set(clone, 'birthday', moment(clone.birthday).format(DateTimeFormatInfo.STANDARD_DATE));
    return clone;
  }
}
