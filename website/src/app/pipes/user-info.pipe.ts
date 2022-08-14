import {Pipe, PipeTransform} from '@angular/core';
import {User} from '../core/models/user.model';
import {Role} from '../core/constants/role.const';

@Pipe({
  name: 'userInfo'
})
export class UserInfoPipe implements PipeTransform {
  transform(user: User, info: 'fullname' | 'username' | 'firstCharName' | 'color' | 'is-doctor'): unknown {
    switch (info) {
      case 'fullname':
        return !!user?.fullName ? user?.fullName : user?.username;
      case 'username':
        return `@${user?.username}`;
      case 'firstCharName':
        return !!user?.lastName
          ? user?.lastName?.charAt(0)?.toUpperCase()
          : user?.username?.charAt(0)?.toUpperCase();
      case 'color':
        return !!user?.color ? user?.color : '#f56a00';
      case 'is-doctor':
        return user?.role === Role.DOCTOR;
    }
  }
}
