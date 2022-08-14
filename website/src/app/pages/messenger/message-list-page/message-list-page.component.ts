import {
  Component,
  Input,
  OnInit,
  Output,
  EventEmitter,
  OnChanges,
  SimpleChanges,
  ViewChild,
  OnDestroy, ChangeDetectionStrategy
} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {GroupService} from '../../../core/services/collectors/group.service';
import {ActivatedRoute, Router} from '@angular/router';
import {User} from '../../../core/models/user.model';
import {StorageService} from '../../../core/services/commons/storage.service';
import * as _ from 'lodash';
import {MessageCreatePageComponent} from '../message-create-page/message-create-page.component';
import {GroupPayload} from '../../../core/models/group/group.payload';
import {NzMessageService} from 'ng-zorro-antd/message';
import {BehaviorSubject, Subject} from 'rxjs';
import {GenerateColorService} from '../../../core/services/commons/generate-color.service';
import {GroupDataSource} from '../../../core/models/datasources/group.datasource';
import {MessageService} from '../../../core/services/ws/message.service';

@Component({
  selector: 'app-message-list-page',
  templateUrl: './message-list-page.component.html',
  styleUrls: ['./message-list-page.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class MessageListPageComponent implements OnInit, OnChanges, OnDestroy {

  @ViewChild('messageCreatePageComponent') messageCreatePageComponent: MessageCreatePageComponent;

  @Input() currentGroupId: number;

  @Output() currentGroupChange = new EventEmitter<Group>();

  public keyword: string;
  public currentUser: User;
  public visibleModalCreateGroup = new BehaviorSubject<boolean>(false);
  public groupDataSource: GroupDataSource;

  private destroy$ = new Subject();

  constructor(private groupService: GroupService,
              private route: ActivatedRoute,
              private router: Router,
              private storageService: StorageService,
              private generateColorService: GenerateColorService,
              private nzMessageService: NzMessageService,
              private messageService: MessageService) {
    this.currentUser = this.storageService.userInfo;

    if (!!this.currentUser) {
      this.groupDataSource = new GroupDataSource(this.groupService, this.messageService, this.generateColorService);
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!changes.currentGroupId && !!this.currentGroupId) {
      if (changes.currentGroupId.currentValue !== changes.currentGroupId.previousValue) {
        this.groupDataSource.setCurrentGroup(this.currentGroupId);
      }
    }
  }

  ngOnInit(): void {
    this.groupDataSource.currentGroup
      .subscribe(group => this.onSelectGroup(group));
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  // region Modal Create Group
  public toggleModalCreateGroup(): void {
    const currentVisible = this.visibleModalCreateGroup.value;
    this.visibleModalCreateGroup.next(!currentVisible);
  }

  public onCreateGroup(): void {
    const groupPayload: GroupPayload = this.messageCreatePageComponent.getModelCreateGroup();

    if (!!groupPayload) {
      this.groupService.createGroup(groupPayload)
        .subscribe(groupCreated => {
          if (!!groupCreated) {
            this.addGroup(groupCreated);
            this.nzMessageService.success('Tạo nhóm thành công');
            this.toggleModalCreateGroup();
          }
        });
    } else {
      this.nzMessageService.warning('Vui lòng chọn thành viên để tạo cuộc hội thoại !');
    }
  }

  // endregion

  public onRefreshListGroup(): void {
    this.groupDataSource.refresh();
    this.groupDataSource.fetchingData(1);
  }

  public onSelectGroup(group: Group): void {
    if (!!group) {
      this.currentGroupId = group.id;
      group.totalNewMessage = 0;
      this.currentGroupChange.emit(group);

      if (group.id !== this.currentGroupId) {
        this.router.navigate(['/messages', group.id])
          .then(() => {
          });
      }
    } else {
      this.router.navigate(['/messages'])
        .then(() => {
        });
    }
  }

  public deleteGroup(group: Group): void {
    if (!!group && group.id) {
      this.groupDataSource.deleteGroup(group.id);
    }
  }

  public addGroup(group: Group): void {
    this.groupDataSource.addGroup(group);
  }

  public updateGroup(group: Group): void {
    this.groupDataSource.updateGroup(group);
  }

  public onSearch(keyword: string): void {
    this.keyword = keyword;
    if (!!this.keyword && this.keyword.length > 0) {
      const oldValue = _.cloneDeep(this.keyword);

      setTimeout(() => {
        if (this.keyword === oldValue) {
          this.groupDataSource.fetchingData(1, this.keyword);
        }
      }, 500);
    } else {
      this.groupDataSource.fetchingData(1);
    }
  }
}
