import {Component, OnDestroy, OnInit, ViewChild} from '@angular/core';
import {NzResizeEvent} from 'ng-zorro-antd/resizable';
import {UserService} from '../../core/services/collectors/user.service';
import {User} from '../../core/models/user.model';
import {StorageService} from '../../core/services/commons/storage.service';
import {KeycloakService} from '../../core/services/auth/keycloak.service';
import {ActivatedRoute} from '@angular/router';
import {Group} from '../../core/models/group/group.model';
import {GroupService} from '../../core/services/collectors/group.service';
import {CachingService} from '../../core/services/commons/caching.service';
import {MessageService} from '../../core/services/ws/message.service';
import {MessageListPageComponent} from '../../pages/messenger/message-list-page/message-list-page.component';
import {filter, retry} from 'rxjs/operators';
import {NzMessageService} from 'ng-zorro-antd/message';
import * as _ from 'lodash';
import {RequestPayload} from '../../core/models/request/request.payload';
import {RequestStatus} from '../../core/constants/request-status.enum';
import {NzModalService} from 'ng-zorro-antd/modal';

@Component({
  selector: 'app-messenger-layout',
  templateUrl: './messenger-layout.component.html',
  styleUrls: ['./messenger-layout.component.scss']
})
export class MessengerLayoutComponent implements OnInit, OnDestroy {

  @ViewChild('messageListPageComponent') messageListPageComponent: MessageListPageComponent;

  public widthSiderRight = 300;
  public widthSiderLeft = 350;
  public idSiderRight = -1;
  public idSiderLeft = -1;
  public currentUser: User;
  public isModalUserInfoVisible = false;
  public currentGroupId: number;
  public currentGroup: Group;
  public collapseInfoTab = false;
  public loading: boolean;

  constructor(private userService: UserService,
              private keycloakService: KeycloakService,
              private storageService: StorageService,
              private groupService: GroupService,
              private messageService: MessageService,
              private cachingService: CachingService,
              private nzMessageService: NzMessageService,
              private route: ActivatedRoute,
              private modal: NzModalService) {
    this.currentUser = this.storageService.userInfo;
    this.route.params
      .subscribe(params => {
        if (!!params) {
          this.currentGroupId = parseInt(params.id, 0);
        }
      });
  }

  ngOnInit(): void {
    if (!!this.currentUser && this.currentGroupId !== undefined) {
      this.messageService.connect();
    }
  }

  ngOnDestroy(): void {
  }

  public onSideResize({width}: NzResizeEvent, isSiderRight: boolean = false): void {
    cancelAnimationFrame(isSiderRight ? this.idSiderRight : this.idSiderLeft);

    if (isSiderRight) {
      this.idSiderRight = requestAnimationFrame(() => {
        this.widthSiderRight = width;
      });
    } else {
      this.idSiderLeft = requestAnimationFrame(() => {
        this.widthSiderLeft = width;
      });
    }
  }

  public onLogout(): void {
    this.userService.logout().subscribe(() => {
      this.keycloakService.logout({
        redirectUri: `${window.location.origin}/auth`,
      });
    });
  }

  public onCurrentGroupChange(group: Group): void {
    if (!!group && (group.id !== this.currentGroupId || !!!this.currentGroup) && !!group.members) {
      this.loading = true;
      this.currentGroupId = group.id;

      const members: Array<User> = new Array<User>();
      this.groupService.getAllMemberOfGroup(group.id)
        .pipe(filter(member => !!member))
        .subscribe(member => members.push(member),
          err => {
          },
          () => {
            group.members = members;
            this.currentGroup = _.cloneDeep(group);
            this.collapseInfoTab = true;
            this.loading = false;
          });
    }
  }

  public onDeleteGroupEvent(group: Group): void {
    this.messageListPageComponent.deleteGroup(group);
  }

  public onAddGroupEvent(group: Group): void {
    this.messageListPageComponent.addGroup(group);
  }

  public onEditGroupEvent(group: Group): void {
    this.messageListPageComponent.updateGroup(group);
  }

  public onRequestJoin(): void {
    this.modal.confirm({
      nzTitle: 'Bạn có chắc muốn tham gia nhóm này ?',
      nzContent: 'Khi nhấn tham gia người dùng trong nhóm có thể thấy các thông tin của bạn',
      nzOkText: 'Tham gia',
      nzCancelText: 'Huỷ',
      nzOnOk: () => {
        if (!!this.currentGroup && !this.currentGroup.isPrivate) {
          const requestPayload = new RequestPayload(
            this.currentGroup.id,
            this.currentUser.userId,
            this.currentUser.userId
          );

          this.groupService.createRequest(requestPayload)
            .subscribe(requestDto => {
              if (requestDto.status === RequestStatus.PENDING) {
                this.nzMessageService.info('Đã gửi yêu cầu tham gia nhóm');
              } else if (requestDto.status === RequestStatus.APPROVE) {
                this.nzMessageService.success('Tham gia nhóm thành công');
                this.currentGroup.isMember = true;
                this.messageListPageComponent.updateGroup(this.currentGroup);
              }
            }, err => this.nzMessageService.error(err));
        } else {
          this.nzMessageService.warning('Bạn không có quyền tham gia vào nhóm này');
        }
      }
    });
  }

  // region User Info Modal
  public onOpenModalUserInfo(): void {
    if (!!this.currentUser) {
      this.isModalUserInfoVisible = true;
    }
  }

  public onCloseModalUserInfo(): void {
    this.isModalUserInfoVisible = false;
  }

  // endregion
}
