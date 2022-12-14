import {
  ChangeDetectionStrategy,
  Component,
  EventEmitter,
  Input,
  OnChanges, OnDestroy,
  OnInit,
  Output,
  SimpleChanges,
  ViewChild
} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {NzModalService} from 'ng-zorro-antd/modal';
import {User} from '../../../core/models/user.model';
import {NzMessageService} from 'ng-zorro-antd/message';
import {GroupService} from '../../../core/services/collectors/group.service';
import {CachingService} from '../../../core/services/commons/caching.service';
import {GroupType} from '../../../core/constants/group-type.const';
import {GroupPayload} from '../../../core/models/group/group.payload';
import {StorageService} from '../../../core/services/commons/storage.service';
import * as _ from 'lodash';
import {SearchUsersComponent} from '../../user/search-users/search-users.component';
import {BehaviorSubject, from} from 'rxjs';
import {MessageEditInfoPageComponent} from '../message-edit-info-page/message-edit-info-page.component';
import {filter} from 'rxjs/operators';
import {GroupRolePipe} from '../../../pipes/group/group-role.pipe';
import {UserInfoPipe} from '../../../pipes/user-info.pipe';
import {RequestPayload} from '../../../core/models/request/request.payload';
import {RequestStatus} from '../../../core/constants/request-status.enum';
import {FileService} from '../../../core/services/collectors/file.service';
import {UploadFileDto} from '../../../core/models/file/upload-file.dto';

@Component({
  selector: 'app-message-info-page',
  templateUrl: './message-info-page.component.html',
  styleUrls: ['./message-info-page.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MessageInfoPageComponent implements OnInit, OnChanges {

  @ViewChild('searchUsersComponent') searchUsersComponent: SearchUsersComponent;
  @ViewChild('messageEditInfoPageComponent') messageEditInfoPageComponent: MessageEditInfoPageComponent;

  @Input() currentGroup: Group;

  @Output() deleteGroupEvent = new EventEmitter<Group>();
  @Output() addGroupEvent = new EventEmitter<Group>();
  @Output() editGroupEvent = new EventEmitter<Group>();

  public collapseMembers: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public collapseImages: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public collapseFiles: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public collapseOptions: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public loading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public visibleModalAddMember: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public visibleModalViewAllMember: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public visibleModalCreateQRCode: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public visibleModalEdit: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public visibleModalViewPendingMembers: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public visibleDrawerFileList: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public currentUser: User;
  public listMemberOld: Array<User> = new Array<User>();
  public groupClone: Group;
  public isListFileDrawer: boolean;

  constructor(private modalService: NzModalService,
              private nzMessageService: NzMessageService,
              private groupService: GroupService,
              private cachingService: CachingService,
              private storageService: StorageService,
              private groupRolePipe: GroupRolePipe,
              private fileService: FileService,
              private userInfoPipe: UserInfoPipe) {
    this.currentUser = this.storageService.userInfo;
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!changes.currentGroup && !!this.currentGroup) {
      if (!this.groupClone || (!!this.groupClone && this.groupClone.id !== this.currentGroup.id)) {
        this.groupClone = _.cloneDeep(this.currentGroup);

        this.collapseMembers.next(false);
        this.collapseImages.next(false);
        this.collapseFiles.next(false);
        this.collapseOptions.next(false);

        this.fileService.getListImageOfGroup(this.currentGroup.id)
          .pipe(filter(image => this.currentGroup.images.length < 9))
          .pipe(filter(image => !this.currentGroup.images.find(iter => iter.nameFile === image.nameFile)))
          .subscribe(image => this.currentGroup.images.push(image));

        this.fileService.getListFileOfGroup(this.currentGroup.id)
          .pipe(filter(file => this.currentGroup.files.length < 5))
          .pipe(filter(file => !this.currentGroup.files.find(iter => iter.nameFile === file.nameFile)))
          .subscribe(file => this.currentGroup.files.push(file));
      }
    }
  }

  // region Modal add member
  public onClickOpenModalAddMember(): void {
    this.listMemberOld = _.clone(this.currentGroup.members);
    this.visibleModalAddMember.next(true);
  }

  public onClickCloseModalAddMember(): void {
    this.visibleModalAddMember.next(false);
  }

  public onAddMembers(): void {
    if (!!this.currentGroup) {
      const users = this.searchUsersComponent.getUserSelected();
      const groupId = this.currentGroup.id;
      const currentUserId = this.currentUser.userId;

      if (this.groupRolePipe.transform(this.currentGroup, 'owner')) {
        this.groupService.addMemberOfGroup(groupId, users)
          .subscribe(result => {
            if (result) {
              this.nzMessageService.success('Th??m th??nh vi??n th??nh c??ng');

              const members = _.cloneDeep(this.currentGroup.members);
              members.push(...users);
              this.currentGroup.members = members;
              this.editGroupEvent.emit(this.currentGroup);
            } else {
              this.nzMessageService.error('Kh??ng th??? th??m th??nh vi??n m???i v??o l??c n??y.');
            }
          });
      } else if (this.userInfoPipe.transform(this.currentUser, 'is-doctor')) {
        from(users)
          .pipe(filter(user => !!user))
          .subscribe(user => {
            const requestPayload = new RequestPayload(groupId, currentUserId, user?.userId);

            this.groupService.createRequest(requestPayload)
              .pipe(filter(requestResult => !!requestResult))
              .subscribe(requestResult => {
                switch (requestResult.status) {
                  case RequestStatus.APPROVE:
                    const members = _.cloneDeep(this.currentGroup.members);
                    members.push(user);
                    this.currentGroup.members = members;

                    this.editGroupEvent.emit(this.currentGroup);
                    this.nzMessageService.success(`???? th??m th??nh vi??n ${user.fullName} th??nh c??ng`);
                    break;
                  case RequestStatus.PENDING:
                    this.nzMessageService.info('Ng?????i d??ng s??? ???????c th??m v??o nh??m khi tr?????ng nh??m ch???p nh???n');
                    break;
                  case RequestStatus.REJECT:
                    this.nzMessageService.error(`Kh??ng ???????c ph??p th??m ${user.fullName} v??o nh??m`);
                    break;
                }
              });
          }, err => this.nzMessageService.error(err));
      }
    } else {
      this.nzMessageService.warning('B???n kh??ng c?? quy???n th??m th??nh vi??n m???i trong nh??m n??y !');
    }

    this.onClickCloseModalAddMember();
  }

  // endregion

  // region Event
  public onClickOpenModalCreateQRCode(): void {
    this.visibleModalCreateQRCode.next(true);
  }

  public onOpenModalEditGroup(): void {
    if (!!this.currentGroup) {
      this.groupClone = _.cloneDeep(this.currentGroup);
      this.visibleModalEdit.next(true);
    }
  }

  public updateGroup(): void {
    if (this.groupRolePipe.transform(this.currentGroup, 'owner')) {
      const groupPayload = this.messageEditInfoPageComponent.getFormValue();

      this.groupService.updateGroup(this.currentGroup.id, groupPayload)
        .pipe(filter(group => !!group))
        .subscribe(groupUpdated => {
          if (!!groupUpdated) {
            this.nzMessageService.success('C???p nh???t th??ng tin th??nh c??ng');

            groupUpdated.color = this.currentGroup.color;
            groupUpdated.members = this.currentGroup.members;
            this.currentGroup = groupUpdated;
            this.editGroupEvent.emit(groupUpdated);
          } else {
            this.nzMessageService.error('L???i c???p nh???t th??ng tin nh??m');
          }

          this.visibleModalEdit.next(false);
        }, err => this.nzMessageService.error(err));
    } else {
      this.nzMessageService.warning('B???n kh??ng c?? quy???n c???p nh???t th??ng tin nh??m n??y');
      this.visibleModalEdit.next(false);
    }
  }

  public onDeleteGroup(): void {
    if (!!this.currentGroup) {
      let message = '';
      const isOwner = this.groupRolePipe.transform(this.currentGroup, 'owner');

      if (this.currentGroup.isGroup) {
        message = `B???n c?? mu???n ${isOwner ? 'xo??' : 'r???i kh???i'} nh??m <b>${this.currentGroup.nameGroup}</b> n??y kh??ng ?`;
      } else {
        message = `B???n c?? mu???n xo?? cu???c trao ?????i v???i <b>${this.currentGroup.nameGroup}</b> n??y kh??ng ?`;
      }

      this.modalService.confirm({
        nzTitle: `C???nh b??o`,
        nzContent: message,
        nzCentered: true,
        nzClosable: false,
        nzAutofocus: 'cancel',
        nzOkText: isOwner ? 'Xo?? cu???c h???i tho???i' : 'R???i cu???c h???i tho???i',
        nzOkDanger: true,
        nzCancelText: 'Hu???',
        nzOnOk: () => isOwner
          ? this.deleteGroup(this.currentGroup.id)
          : this.outGroup(this.currentGroup.id)
      });
    }
  }

  public onCreateMessenger(userId: string): void {
    const groupPayload: GroupPayload = {
      type: GroupType.ONE,
      private: true,
      users: [userId],
      description: null,
      nameGroup: null,
    };

    this.groupService.createGroup(groupPayload)
      .subscribe(group => {
          if (!!group) {
            this.addGroupEvent.emit(group);
            this.visibleModalViewAllMember.next(false);
          }
        }, err => this.visibleModalViewAllMember.next(false),
        () => this.visibleModalViewAllMember.next(false));
  }

  public onClickDeleteMember(member: User): void {
    const isOwner = this.groupRolePipe.transform(this.currentGroup, 'owner');
    if (!!member && isOwner) {
      this.modalService.warning({
        nzTitle: 'Xo?? th??nh vi??n',
        nzContent: `B???n c?? mu???n xo?? <span class="font-bold">${member?.fullName}</span> ra kh???i nh??m kh??ng ?`,
        nzCentered: true,
        nzAutofocus: 'cancel',
        nzCancelText: 'Hu???',
        nzOkText: 'Xo?? th??nh vi??n',
        nzOkDanger: true,
        nzOnOk: () => {
          if (!!this.currentGroup && !!member) {
            this.groupService.deleteMemberOfGroup(this.currentGroup.id, member.userId)
              .subscribe(result => {
                if (result) {
                  this.nzMessageService.success('???? xo?? th??nh vi??n kh???i nh??m');

                  const members = _.cloneDeep(this.currentGroup.members);
                  this.currentGroup.members = members.filter(iter => iter.userId !== member.userId);
                  this.editGroupEvent.emit(this.currentGroup);
                } else {
                  this.nzMessageService.error('Kh??ng th??? xo?? th??nh vi??n kh???i nh??m l??c n??y');
                }
              }, err => this.nzMessageService.error(err));
          }
        }
      });
    }
  }

  public onOpenFileList(): void {
    this.isListFileDrawer = true;
    this.visibleDrawerFileList.next(true);
  }

  public onOpenImageList(): void {
    this.isListFileDrawer = false;
    this.visibleDrawerFileList.next(true);
  }

  public onDeleteUserFromListFull(userId: string): void {
    const members = _.cloneDeep(this.currentGroup.members);
    this.currentGroup.members = members.filter(member => member.userId !== userId);
    this.editGroupEvent.emit(this.currentGroup);
  }

  // endregion

  private deleteGroup(groupId: number): void {
    const messId = this.nzMessageService.loading(
      '??ang x??a cu???c h???i tho???i c???a b???n ...',
      {nzDuration: 0}
    ).messageId;

    this.groupService.deleteGroup(groupId).subscribe(
      (result) => {
        this.nzMessageService.remove(messId);

        if (result) {
          this.nzMessageService.success('???? x??a cu???c h???i tho???i.');
          this.deleteGroupEvent.emit(this.currentGroup);
        } else {
          this.nzMessageService.error(
            'Kh??ng th??? x??a cu???c h???i tho???i v??o l??c n??y. Vui l??ng th??? l???i sau'
          );
        }
      },
      (error) => {
        this.nzMessageService.remove(messId);
        this.nzMessageService.error(error);
      },
      () => this.nzMessageService.remove(messId)
    );
  }

  private outGroup(groupId: number): void {
    const messId = this.nzMessageService.loading(
      '??ang r???i kh???i cu???c h???i tho???i c???a n??y ...',
      {nzDuration: 0}
    ).messageId;

    this.groupService.memberOutGroup(groupId).subscribe(
      (result) => {
        this.nzMessageService.remove(messId);

        if (result) {
          this.currentGroup.isMember = false;
          this.nzMessageService.success('???? r???i kh???i cu???c h???i tho???i.');
          this.editGroupEvent.emit(this.currentGroup);
          if (this.currentGroup.isPrivate) {
            this.deleteGroupEvent.emit(this.currentGroup);
          }
        } else {
          this.nzMessageService.error(
            'Kh??ng th??? r???i kh???i cu???c h???i tho???i v??o l??c n??y. Vui l??ng th??? l???i sau'
          );
        }
      },
      (error) => {
        this.nzMessageService.remove(messId);
        this.nzMessageService.error(error);
      },
      () => this.nzMessageService.remove(messId)
    );
  }
}
