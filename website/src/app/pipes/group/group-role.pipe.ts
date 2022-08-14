import {Pipe, PipeTransform} from '@angular/core';
import {User} from '../../core/models/user.model';
import {Group} from '../../core/models/group/group.model';
import {StorageService} from '../../core/services/commons/storage.service';

@Pipe({
  name: 'groupRole'
})
export class GroupRolePipe implements PipeTransform {

  private currentUser: User;

  constructor(private storageService: StorageService) {
    this.currentUser = this.storageService.userInfo;
  }


  transform(group: Group, role: 'owner' | 'member', user?: User): boolean {
    switch (role) {
      case 'member':
        return this.isMember(group, user);
      case 'owner':
        return this.isOwner(group, user);
      default:
        return false;
    }
  }

  private isOwner(group: Group, user?: User): boolean {
    if (!!user && !!this.currentUser) {
      this.currentUser = this.storageService.userInfo;
    }

    return !!user ? group?.owner === user.userId : group?.owner === this.currentUser.userId;
  }

  private isMember(group: Group, user?: User): boolean {
    if (!!user && !!this.currentUser) {
      this.currentUser = this.storageService.userInfo;
    }

    return group?.isMember;
  }
}
