import {Component, Input, OnInit, ViewChild} from '@angular/core';
import {User} from '../../../core/models/user.model';
import {Group} from '../../../core/models/group/group.model';
import {FormControl, FormGroup} from '@angular/forms';
import {GroupPayload} from '../../../core/models/group/group.payload';
import * as _ from 'lodash';
import {GroupType} from '../../../core/constants/group-type.const';
import {SearchUsersComponent} from '../../user/search-users/search-users.component';
import {AsyncSubject} from 'rxjs';

@Component({
  selector: 'app-message-create-page',
  templateUrl: './message-create-page.component.html',
  styleUrls: ['./message-create-page.component.scss']
})
export class MessageCreatePageComponent implements OnInit {

  @ViewChild('searchUsersComponent') searchUsersComponent: SearchUsersComponent;

  @Input() currentUser: User;

  public formCreateGroup: FormGroup;
  public totolSelected = new AsyncSubject<number>();

  constructor() {
    this.formCreateGroup = this.createFormGroup();
    this.formCreateGroup.disable();
  }

  ngOnInit(): void {
  }

  public onUserSelectedChange(totalSelected: number): void {
    this.totolSelected.next(totalSelected);
    if (totalSelected >= 2) {
      this.formCreateGroup.enable();
    } else {
      this.formCreateGroup.disable();
    }
  }

  public getModelCreateGroup(): GroupPayload {
    const users: Array<User> = this.searchUsersComponent.getUserSelected();

    if (!!users && users.length >= 1) {
      const formValue = this.formCreateGroup.getRawValue();

      const nameGroupDefault = users.map(user => user.firstName).join(', ');

      const nameGroup = _.get(formValue, 'nameGroup', nameGroupDefault);
      const isPublicGroup: boolean = _.get(formValue, 'publicGroup', false);

      return {
        nameGroup: nameGroup.trim().length > 0 ? nameGroup : nameGroupDefault,
        description: '',
        private: !isPublicGroup,
        type: users.length === 1 ? GroupType.ONE : GroupType.MANY,
        users: _.uniq(users.map(user => user.userId).filter(userId => userId !== this.currentUser.userId))
      };
    }

    return null;
  }

  private createFormGroup(group?: Group): FormGroup {
    const formGroup: FormGroup = new FormGroup({
      nameGroup: new FormControl(''),
      publicGroup: new FormControl(false)
    });

    if (!!group) {
      formGroup.patchValue(group);
    }

    return formGroup;
  }
}
