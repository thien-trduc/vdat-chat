import {
  Component,
  Input,
  OnInit,
  Output,
  EventEmitter,
  OnDestroy,
  ViewChild,
  ElementRef,
  OnChanges, SimpleChanges
} from '@angular/core';
import {User} from '../../../core/models/user.model';
import {UserService} from '../../../core/services/collectors/user.service';
import {GenerateColorService} from '../../../core/services/commons/generate-color.service';
import {StorageService} from '../../../core/services/commons/storage.service';
import * as _ from 'lodash';
import {Subject} from 'rxjs';
import {takeUntil} from 'rxjs/operators';
import {UserDataSource} from '../../../core/models/datasources/user.datasource';
import {SelectItem} from '../../../core/models/select-item';

@Component({
  selector: 'app-search-users',
  templateUrl: './search-users.component.html',
  styleUrls: ['./search-users.component.scss']
})
export class SearchUsersComponent implements OnInit, OnChanges, OnDestroy {
  @ViewChild('searchInputElement') searchInputElement: ElementRef;

  @Input() addCurrentUser: boolean;
  @Input() ignoreUsers: Array<User> = new Array<User>();
  @Output() totalSelected = new EventEmitter<number>();

  public users: Array<User>;
  public loading: boolean;
  public currentUser: User;
  public userDataSource: UserDataSource;
  public keyword: string;

  private destroy$ = new Subject();

  constructor(private userService: UserService,
              private generateColorService: GenerateColorService,
              private storageService: StorageService) {
    this.users = new Array<User>();
    this.currentUser = this.storageService.userInfo;
    this.userDataSource = new UserDataSource(this.userService, this.generateColorService,
      this.currentUser, this.ignoreUsers);
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes) {
      if (changes.ignoreUsers) {
        this.userDataSource.setIgnoreUser(this.ignoreUsers);
      }
    }
  }

  ngOnInit(): void {
    this.userDataSource
      .completed()
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => {
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  public onSearchUsers(keyword: string): void {
    this.keyword = keyword;
    if (!!this.keyword && this.keyword.length > 0) {
      const oldValue = _.cloneDeep(this.keyword);

      setTimeout(() => {
        if (this.keyword === oldValue) {
          this.userDataSource.fetchingData(1, this.keyword);
        }
      }, 500);
    } else {
      this.userDataSource.fetchingData(1);
    }
  }

  public onToggleSelectUser(item: SelectItem<User>): void {
    this.userDataSource.toggleItemSelected(item);
    this.totalSelected.emit(this.userDataSource.selecteds.length);
  }

  public getUserSelected(): Array<User> {
    return this.userDataSource.selecteds
      .map(iter => iter.data);
  }
}
