import {BaseDataSource} from './base.datasource';
import * as _ from 'lodash';
import {GenerateColorService} from '../../services/commons/generate-color.service';
import {CollectionViewer} from '@angular/cdk/collections';
import {filter, map, takeUntil} from 'rxjs/operators';
import {Group} from '../group/group.model';
import {GroupService} from '../../services/collectors/group.service';
import {User} from '../user.model';
import {BehaviorSubject, Subject} from 'rxjs';
import {MessageService} from '../../services/ws/message.service';
import {Message} from '../message/message.model';

export class GroupDataSource extends BaseDataSource<Group> {
  public currentGroup: Subject<Group>;
  private currentGroupId: number;
  private keyword: string;
  private selected: Group;

  constructor(private groupService: GroupService,
              private messageService?: MessageService,
              private generateColorService?: GenerateColorService) {
    super();
    this.currentGroup = new Subject<Group>();

    this.messageService.newMessage$
      .pipe(filter(messageResponse => messageResponse.groupId > 0))
      .subscribe(messageResponse => {
        const group: Group = this.cachedData.find(iter => iter.id === messageResponse.groupId);

        if (!!group) {
          group.lastMessage = messageResponse.message;
          group.totalNewMessage += 1;
          group.updatedAt = new Date();

          this.cachedData = _.orderBy(this.cachedData, 'updatedAt', 'desc');
          this.dataStream.next(this.cachedData);
        }
      });

    this.messageService.deleteMessage$
      .pipe(filter(messageResponse => messageResponse.groupId > 0))
      .subscribe(messageResponse => {
        const group: Group = this.cachedData.find(iter => iter.id === messageResponse.groupId);
        const lastMessage: Message = _.cloneDeep(_.get(group, 'lastMessage', null));

        if (!!group && !!lastMessage && lastMessage.id === messageResponse.messageId) {
          lastMessage.message = '';
          lastMessage.deletedAt = new Date();

          group.updatedAt = new Date();
          group.lastMessage = lastMessage;

          this.cachedData = _.orderBy(this.cachedData, 'updatedAt', 'desc');
          this.dataStream.next(this.cachedData);
        }
      });
  }

  public setCurrentGroup(groupId: number): void {
    if (this.currentGroupId !== groupId) {
      this.currentGroupId = groupId;
      const groupFound = this.cachedData.find(group => group.id === groupId);
      if (!!groupFound) {
        this.currentGroup.next(groupFound);
        this.selected = groupFound;
      }
    }
  }

  public deleteGroup(groupId: number): void {
    const index = this.cachedData.findIndex(group => group.id === groupId);
    this.cachedData.splice(index, 1);
    this.dataStream.next(this.cachedData);

    if (!!this.cachedData && this.cachedData.length > 0) {
      const currentGroup: Group = this.cachedData[0];
      this.setCurrentGroup(currentGroup.id);
    }
  }

  public addGroup(group: Group): void {
    const groupFound: Group = this.cachedData.find(iter => iter.id === group.id);

    if (!!!groupFound) {
      this.cachedData.splice(0, 0, group);
      this.dataStream.next(this.cachedData);
      this.setCurrentGroup(group.id);
    }
  }

  public updateGroup(group: Group): void {
    const groupIndex = this.cachedData.findIndex(iter => iter.id === group.id);

    if (groupIndex !== -1) {
      this.cachedData.splice(groupIndex, 1, group);
      this.dataStream.next(this.cachedData);
    }
  }

  protected setup(collectionViewer: CollectionViewer): void {
    this.fetchingData(1);
    collectionViewer.viewChange
      .pipe(
        takeUntil(this.complete$),
        takeUntil(this.disconnect$))
      .subscribe(range => {
        const endPage = this.getPageForIndex(range.end);
        this.fetchingData(endPage + 1, this.keyword);
      });
  }

  private getPageForIndex(index: number): number {
    return Math.floor(index / this.pageSize);
  }

  public fetchingData(page: number, keyword?: string): void {
    if (this.keyword !== keyword) {
      this.keyword = keyword;
      this.refresh();
    } else if (this.fetchedPages.has(page)) {
      return;
    }

    this.toogleLoading(true);
    this.fetchedPages.add(page);

    this.groupService.getAllGroup(page, this.pageSize, this.keyword)
      .subscribe(group => {
        this.cachedData.splice(page * this.pageSize, this.pageSize, group);
        this.cachedData = _.uniqBy(this.cachedData, 'id');
        this.dataStream.next(this.cachedData);

        if (this.currentGroupId !== undefined) {
          this.setCurrentGroup(this.currentGroupId);
        } else if (page === 1 && !!this.cachedData && this.cachedData.length > 0) {
          const currentGroup: Group = this.cachedData[0];
          this.setCurrentGroup(currentGroup.id);
        }

        this.toogleLoading(false);
      });
  }
}
