import {AfterViewInit, Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {BehaviorSubject} from 'rxjs';
import {RequestDto} from '../../../core/models/request/request.dto';
import {GroupService} from '../../../core/services/collectors/group.service';
import {NzNotificationService} from 'ng-zorro-antd/notification';

@Component({
  selector: 'app-message-list-pending-members-page',
  templateUrl: './message-list-pending-members-page.component.html',
  styleUrls: ['./message-list-pending-members-page.component.scss']
})
export class MessageListPendingMembersPageComponent implements OnInit, OnChanges {

  @Input() currentGroup: Group;

  public loading = new BehaviorSubject<boolean>(false);
  public listMemberPending: Array<RequestDto>;

  constructor(private groupService: GroupService,
              private notificationService: NzNotificationService) {
    this.listMemberPending = new Array<RequestDto>();
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!this.currentGroup) {
      this.loading.next(true);
      this.groupService.getAllRequestInGroup(this.currentGroup.id)
        .subscribe(request => {
            this.listMemberPending.push(request);
            this.loading.next(false);
          }, err => this.loading.next(false),
          () => this.loading.next(false));
    }
  }

  // region Event
  public onApprove(request: RequestDto): void {
    if (!!request) {
      this.groupService.approveRequest(request.requestId)
        .subscribe(result => {
          if (result) {
            this.listMemberPending = this.listMemberPending.filter(iter => iter.requestId !== request.requestId);

            this.notificationService.success('Duyệt thành công',
              `Đã thêm người dùng ${request.userInvited.fullName} vào nhóm`,
              {nzAnimate: true});
          } else {
            this.notificationService.error('Duyệt thất bại',
              `Không thể thêm người dùng ${request.userInvited.fullName} vào nhóm lúc này. Vui lòng thử lại sau`,
              {nzAnimate: true});
          }
        });
    }
  }

  public onReject(request: RequestDto): void {
    if (!!request) {
      this.groupService.rejectRequest(request.requestId)
        .subscribe(result => {
          if (result) {
            this.listMemberPending = this.listMemberPending.filter(iter => iter.requestId !== request.requestId);

            this.notificationService.success('Từ chối thành công',
              `Đã từ chối người dùng ${request.userInvited.fullName} vào nhóm`,
              {nzAnimate: true});
          } else {
            this.notificationService.error('Từ chối thất bại',
              `Không thể từ chối người dùng ${request.userInvited.fullName} vào nhóm lúc này. Vui lòng thử lại sau`,
              {nzAnimate: true});
          }
        });
    }
  }

  // endregion

  public onRefresh(): void {
    if (!!this.currentGroup) {
      this.loading.next(true);
      this.listMemberPending = new Array<RequestDto>();
      this.groupService.getAllRequestInGroup(this.currentGroup.id)
        .subscribe(request => {
            this.listMemberPending.push(request);
            this.loading.next(false);
          }, err => this.loading.next(false),
          () => this.loading.next(false));
    }
  }
}
