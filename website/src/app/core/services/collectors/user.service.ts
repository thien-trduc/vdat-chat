import {Injectable} from '@angular/core';
import {ApiService} from '../commons/api.service';
import {User} from '../../models/user.model';
import {Observable} from 'rxjs';
import {environment} from '../../../../environments/environment';
import {CachingService} from '../commons/caching.service';
import {StorageService} from '../commons/storage.service';
import {KeycloakService} from '../auth/keycloak.service';
import * as _ from 'lodash';
import {GenerateColorService} from '../commons/generate-color.service';
import {filter, map} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class UserService {

  private readonly apiEndpoint = environment.api.users;

  constructor(private apiService: ApiService,
              private cachingService: CachingService,
              private storageService: StorageService,
              private keycloakService: KeycloakService,
              private generateColorService: GenerateColorService) {
  }

  public getUserInfo(): Observable<User> {
    return new Observable<User>(observer => {
      const url = `${this.apiEndpoint}/info`;

      this.cachingService.get<User>(url)
        .subscribe((user) => {
          if (!!user) {
            observer.next(user);
          }
        });

      this.apiService.get(url).subscribe(
        (res) => {
          const body = _.get(res, 'body', []);
          const user = User.fromJson(body, this.generateColorService);
          const userFromKeycloak = this.keycloakService.userInfo;

          if (userFromKeycloak) {
            user.firstName = _.get(userFromKeycloak, 'given_name', '');
            user.lastName = _.get(userFromKeycloak, 'family_name', '');
            user.fullName = _.get(userFromKeycloak, 'name', '');
            user.username = _.get(userFromKeycloak, 'preferred_username', '');
          }

          this.cachingService.save<User>(url, user).subscribe(() => {
          });
          observer.next(user);
        },
        (err) => {
          observer.error(err);
        },
        () => observer.complete()
      );
    });
  }

  public logout(): Observable<any> {
    return new Observable<any>((observer) => {
      const user: User = this.storageService.userInfo;
      const url = `${this.apiEndpoint}/online?socketId=${user.socketId}&hostName=${user.hostName}`;

      this.apiService.delete(url).subscribe(
        (res) => observer.next(true),
        (err) => {
          observer.error(err);
        },
        () => observer.complete()
      );
    });
  }

  public findUserByKeyword(
    keyword: string,
    page?: number,
    pageSize?: number
  ): Observable<User> {
    return new Observable<User>((observer) => {
      const currentUser = this.storageService.userInfo;
      let url = `${this.apiEndpoint}?keyword=${keyword}`;

      if (!!page) {
        url += `&page=${page}`;

        if (!!page && !!pageSize) {
          url += `&pageSize=${pageSize}`;
        }
      }

      this.cachingService.get<Array<User>>(url)
        .pipe(filter(users => !!users))
        .pipe(map(users => users.filter(user => user.userId !== currentUser.userId)))
        .subscribe(users => {
          if (!!users) {
            users.forEach(user => observer.next(user));
          }
        });

      this.apiService.get(url)
        .pipe(map<User, Array<User>>(response => _.get(response, 'body', [])))
        .pipe(map(users => users.map(value => User.fromJson(value))))
        .pipe(map(users => users.filter(user => user.userId !== currentUser.userId)))
        .subscribe(users => {
            this.cachingService.save<Array<User>>(url, users).subscribe(() => {
            });

            users.forEach(user => observer.next(user));
          }, err => {
            observer.error(err);
          },
          () => observer.complete()
        );
    });
  }
}
