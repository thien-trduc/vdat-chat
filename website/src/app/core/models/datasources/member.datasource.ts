import {User} from '../user.model';
import * as _ from 'lodash';
import {GenerateColorService} from '../../services/commons/generate-color.service';
import {CollectionViewer} from '@angular/cdk/collections';
import {filter, takeUntil} from 'rxjs/operators';
import {GroupService} from '../../services/collectors/group.service';
import {Role} from '../../constants/role.const';
import {BaseDataSource} from './base.datasource';
import {BehaviorSubject} from 'rxjs';

export class MemberDataSource extends BaseDataSource<User> {
  private listAllMember: Array<User> = new Array<User>();
  public isFilterAll: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
  public isFilterDoctorOnly: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  public isFilterPatientOnly: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  constructor(private groupService: GroupService,
              private generateColorService: GenerateColorService,
              private groupId: number,
              private keyword?: string) {
    super();
  }

  // region Filter Data
  public filterData(): void {
    if (this.isFilterAll.value) {
      this.filterAll();
    } else if (this.isFilterDoctorOnly.value) {
      this.filterDoctorOnly();
    } else {
      this.filterPatientOnly();
    }
  }

  public filterAll(): void {
    this.isFilterAll.next(true);
    this.isFilterDoctorOnly.next(false);
    this.isFilterPatientOnly.next(false);

    this.cachedData = _.cloneDeep(this.listAllMember);
    this.dataStream.next(this.cachedData);
    console.log(this.cachedData);
    this.toogleLoading(false);
  }

  public filterDoctorOnly(): void {
    this.isFilterAll.next(false);
    this.isFilterDoctorOnly.next(true);
    this.isFilterPatientOnly.next(false);

    this.cachedData = this.listAllMember.filter(member => member.role === Role.DOCTOR);
    this.dataStream.next(this.cachedData);
    this.toogleLoading(false);
  }

  public filterPatientOnly(): void {
    this.isFilterAll.next(false);
    this.isFilterDoctorOnly.next(false);
    this.isFilterPatientOnly.next(true);

    this.cachedData = this.listAllMember.filter(member => member.role === Role.PATIENT);
    this.dataStream.next(this.cachedData);
    this.toogleLoading(false);
  }

  // endregion

  public deleteUser(userId: string): void {
    const index = this.listAllMember.findIndex(user => user.userId === userId);
    if (index !== -1) {
      this.listAllMember.splice(index, 1);
      this.filterData();
    }
  }

  public refreshList(): void {
    this.refresh();
    this.listAllMember = new Array<User>();
    this.fetchingData(1);
  }

  protected setup(collectionViewer: CollectionViewer): void {
    this.fetchingData(1);
    collectionViewer.viewChange
      .pipe(
        takeUntil(this.complete$),
        takeUntil(this.disconnect$))
      .subscribe(range => {
        const endPage = this.getPageForIndex(range.end);
        this.fetchingData(endPage + 1);
      });
  }

  private getPageForIndex(index: number): number {
    return Math.floor(index / this.pageSize);
  }

  private fetchingData(page: number): void {
    if (this.fetchedPages.has(page)) {
      return;
    }
    this.fetchedPages.add(page);

    this.toogleLoading(true);
    this.groupService.getAllMemberOfGroup(this.groupId, page, this.pageSize, this.keyword)
      .pipe(filter(user => !this.listAllMember.find(iter => iter.userId === user.userId)))
      .subscribe(member => {
          this.listAllMember.splice(page * this.pageSize, this.pageSize, member);

          this.filterData();

          this.toogleLoading(false);
        }, () => {
          this.toogleLoading(false);
        },
        () => this.toogleLoading(false));
  }
}
