import {User} from '../user.model';
import {UserService} from '../../services/collectors/user.service';
import * as _ from 'lodash';
import {GenerateColorService} from '../../services/commons/generate-color.service';
import {CollectionViewer} from '@angular/cdk/collections';
import {distinct, filter, map, takeUntil} from 'rxjs/operators';
import {SelectItem} from '../select-item';
import {BaseDataSource} from './base.datasource';

export class UserDataSource extends BaseDataSource<SelectItem<User>> {
  private keyword: string;

  constructor(private userService: UserService,
              private generateColorService: GenerateColorService,
              private currentUser: User,
              private ignoreUsers?: Array<User>) {
    super();
  }

  public setIgnoreUser(ignoreUsers: Array<User>): void {
    this.ignoreUsers = ignoreUsers;
    this.cachedData.filter(iter => !!!ignoreUsers.find(user => user.userId === iter.data.userId));
    this.cachedData = _.uniqBy(this.cachedData, 'data.userId');
    this.dataStream.next(this.cachedData);
  }

  public toggleItemSelected(item: SelectItem<User>): void {
    const itemCache = this.cachedData.find(iter => iter.data.userId === item.data.userId);
    const itemFound = this.selecteds.find(iter => iter.data.userId === item.data.userId);

    if (!!itemCache) {
      item.toggleSelect();
      itemCache.toggleSelect();

      if (!!itemFound && !itemFound.checked) {
        this.selecteds = this.selecteds.filter(iter => iter.data.userId !== item.data.userId);
      } else {
        this.selecteds.push(item);
      }
    }
  }

  public popSelected(): void {
    if (this.selecteds.length > 0 && !this.keyword) {
      const lastItemSelected = this.selecteds.pop();
      const itemCache = this.cachedData.find(iter => iter.data.userId === lastItemSelected.data.userId);
      itemCache.unSelect();
    }
  }

  protected setup(collectionViewer: CollectionViewer): void {
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

    this.userService.findUserByKeyword(this.keyword, page, this.pageSize)
      .pipe(filter(user => user.userId !== this.currentUser.userId))
      .pipe(filter(user => !!!this.ignoreUsers.find(iter => iter.userId === user.userId)))
      .pipe(map(user => {
        user.color = this.generateColorService.generate(user.userId);
        return user;
      }))
      .pipe(map(user => {
        const item = new SelectItem<User>(user);
        const itemSelectedFound = this.selecteds.find(iter => iter.data.userId === user.userId);
        item.checked = !!itemSelectedFound ? itemSelectedFound.checked : false;
        return item;
      }))
      .subscribe(user => {
          this.cachedData.splice(page * this.pageSize, this.pageSize, user);
          this.cachedData = _.uniqBy(this.cachedData, 'data.userId');
          this.dataStream.next(this.cachedData);
          this.toogleLoading(false);
        }, err => this.toogleLoading(false),
        () => this.toogleLoading(false));
  }
}
