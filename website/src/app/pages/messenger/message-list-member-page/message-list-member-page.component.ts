import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  OnInit,
  Output,
  SimpleChanges,
  ViewChild
} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {User} from '../../../core/models/user.model';
import {Subject} from 'rxjs';
import {GenerateColorService} from '../../../core/services/commons/generate-color.service';
import {filter, takeUntil} from 'rxjs/operators';
import {GroupService} from '../../../core/services/collectors/group.service';
import {MemberDataSource} from '../../../core/models/datasources/member.datasource';
import {NzMessageService} from 'ng-zorro-antd/message';
import {GroupRolePipe} from '../../../pipes/group/group-role.pipe';

@Component({
  selector: 'app-message-list-member-page',
  templateUrl: './message-list-member-page.component.html',
  styleUrls: ['./message-list-member-page.component.scss']
})
export class MessageListMemberPageComponent implements OnInit, OnChanges, OnDestroy {
  @ViewChild('searchInputElement') searchInputElement: ElementRef;

  @Input() currentGroup: Group;
  @Input() currentUser: User;

  @Output() deleteUserEvent = new EventEmitter<string>();
  @Output() createMessageEvent = new EventEmitter<string>();

  public users: Array<User>;
  public loading: boolean;
  public memberDataSource: MemberDataSource;

  private destroy$ = new Subject();

  constructor(private groupService: GroupService,
              private generateColorService: GenerateColorService,
              private nzMessageService: NzMessageService,
              private groupRolePipe: GroupRolePipe) {
    this.users = new Array<User>();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!this.currentGroup) {
      this.memberDataSource = new MemberDataSource(
        this.groupService,
        this.generateColorService,
        this.currentGroup.id);
    }
  }

  ngOnInit(): void {
    this.memberDataSource
      .completed()
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => {
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // region Event
  public onDeleteUser(userId: string): void {
    if (this.groupRolePipe.transform(this.currentGroup, 'owner')) {
      if (!!this.currentGroup) {
        this.groupService.deleteMemberOfGroup(this.currentGroup.id, userId)
          .subscribe(result => {
            if (result) {
              this.nzMessageService.success('Đã xoá thành viên khỏi nhóm');
              this.deleteUserEvent.emit(userId);
              this.memberDataSource.deleteUser(userId);
            } else {
              this.nzMessageService.error('Không thể xoá thành viên khỏi nhóm lúc này');
            }
          }, err => this.nzMessageService.error(err));
      }
    }
  }

  public onCreateMessage(userId: string): void {
    this.createMessageEvent.emit(userId);
  }
  // endregion
}
