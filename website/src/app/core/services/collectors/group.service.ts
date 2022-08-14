import {Injectable} from '@angular/core';
import {environment} from '../../../../environments/environment';
import {ApiService} from '../commons/api.service';
import {CachingService} from '../commons/caching.service';
import {GenerateColorService} from '../commons/generate-color.service';
import {from, Observable} from 'rxjs';
import {Group} from '../../models/group/group.model';
import * as _ from 'lodash';
import {GroupPayload} from '../../models/group/group.payload';
import {User} from '../../models/user.model';
import {filter, finalize, map} from 'rxjs/operators';
import {RequestPayload} from '../../models/request/request.payload';
import {RequestDto} from '../../models/request/request.dto';
import {RequestStatus} from "../../constants/request-status.enum";

@Injectable({
  providedIn: 'root'
})
export class GroupService {

  private readonly apiEndpoint = environment.api.groups;
  private readonly apiRequestEndpoint = environment.api.request;

  constructor(private apiService: ApiService,
              private cachingService: CachingService,
              private generateColorService: GenerateColorService) {
  }

  // region Group
  public getAllGroup(page?: number, pageSize?: number, keyword?: string): Observable<Group> {
    return new Observable<Group>(observer => {
      let url = `${this.apiEndpoint}`;

      if (!!keyword && keyword.length > 0) {
        url += `?keyword=${keyword}`;
      }

      if (!!page) {
        url += `${!!keyword ? '&' : '?'}page=${page}`;

        if (!!page && !!pageSize) {
          url += `&pageSize=${pageSize}`;
        }
      }

      this.cachingService
        .get<Array<Group>>(url)
        .pipe(filter(groups => !!groups))
        .subscribe(groups => groups.forEach(group => observer.next(group)));

      this.apiService.get(url)
        .pipe(filter(response => !!response))
        .pipe(map<Group, Array<Group>>(response => _.get(response, 'body', [])))
        .pipe(filter(groups => !!groups))
        .pipe(map(groups => groups.map(group => Group.fromJson(group, this.generateColorService))))
        .pipe(map(groups => _.uniqBy(groups, 'id')))
        .pipe(map(groups => _.orderBy(groups, 'updatedAt', 'desc')))
        .pipe(filter(groups => !!groups))
        .subscribe(groups => {
            groups.forEach(group => observer.next(group));

            this.cachingService.save(url, groups).subscribe((result) => {
            });
          }, err => observer.error(err),
          () => observer.complete()
        );
    });
  }

  public createGroup(groupPayload: GroupPayload): Observable<Group> {
    return new Observable<any>((observer) => {
      const url = `${this.apiEndpoint}`;

      this.apiService.post(url, groupPayload).subscribe(
        (res) => {
          const body = _.get(res, 'body', {});
          if (_.isArray(body) && body.length > 0) {
            observer.next(Group.fromJson(body[0]));
          } else {
            observer.next(Group.fromJson(body));
          }
        },
        (error) => observer.error(error),
        () => observer.complete()
      );
    });
  }

  public updateGroup(groupId: number, group: GroupPayload): Observable<Group> {
    return new Observable<Group>((observer) => {
      const url = `${this.apiEndpoint}/${groupId}`;

      this.apiService.put(url, group).subscribe(
        (res) => {
          const body = _.get(res, 'body', {});
          observer.next(Group.fromJson(body));
        },
        (error) => observer.error(error),
        () => observer.complete()
      );
    });
  }

  public deleteGroup(groupId: number): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      const url = `${this.apiEndpoint}/${groupId}`;

      this.apiService.delete(url).subscribe(
        (res) => {
          const body = _.get(res, 'body', {});
          const result = _.get(body, 'result', false);
          observer.next(result);
        },
        (error) => observer.error(error),
        () => observer.complete()
      );
    });
  }

  // endregion

  // region Group Member
  public addMemberOfGroup(groupId: number, users: Array<User>): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      const url = `${this.apiEndpoint}/${groupId}/members`;
      const userIds = users
        .filter((user) => !!user.userId)
        .map((user) => user.userId);

      this.apiService.patch(url, {users: userIds})
        .pipe(map(response => _.get(response, 'body', {})))
        .pipe(map(body => _.get(body, 'result', false)))
        .subscribe(result => {
            observer.next(!!result);
          },
          (error) => observer.error(error),
          () => observer.complete()
        );
    });
  }

  public getAllMemberOfGroup(groupId: number, page?: number, pageSize?: number, keyword?: string): Observable<User> {
    return new Observable<User>((observer) => {
      const url = `${this.apiEndpoint}/${groupId}/members`;

      this.cachingService
        .get<Array<User>>(url)
        .pipe(filter(users => !!users && users.length > 0))
        .subscribe(users => users.forEach(user => observer.next(user)));

      this.apiService.get(url)
        .pipe(
          map<Array<any>, Array<any>>(response => _.get(response, 'body', [])),
          filter(members => !!members),
          map(users => users.map(user => User.fromJson(user, this.generateColorService)))
        )
        .subscribe(members => {
            this.cachingService.save(url, members).subscribe((result) => {
            });

            members.forEach(member => observer.next(member));
          },
          (error) => observer.error(error),
          () => observer.complete()
        );
    });
  }

  public deleteMemberOfGroup(
    groupId: number,
    userId: string
  ): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      const url = `${this.apiEndpoint}/${groupId}/members/${userId}`;

      this.apiService.delete(url).subscribe(
        (res) => {
          const body = _.get(res, 'body', []);
          const result = _.get(body, 'result', false);
          observer.next(!!result);
        },
        (error) => observer.error(error),
        () => observer.complete()
      );
    });
  }

  public memberOutGroup(groupId: number): Observable<boolean> {
    return new Observable<boolean>((observer) => {
      const url = `${this.apiEndpoint}/${groupId}/members`;

      this.apiService.delete(url).subscribe(
        (res) => {
          const body = _.get(res, 'body', []);
          const result = _.get(body, 'result', false);
          observer.next(!!result);
        },
        (error) => observer.error(error),
        () => observer.complete()
      );
    });
  }

  // endregion

  // region Group Request
  public getAllRequestInGroup(groupId: number): Observable<RequestDto> {
    return new Observable<RequestDto>(observer => {
      this.apiService.get(`${this.apiRequestEndpoint}/request/${groupId}`)
        .pipe(filter(response => !!response))
        .pipe(map(response => _.get(response, 'body', [])))
        .pipe(filter(arr => arr.length > 0))
        .subscribe(requests => {
          from(requests)
            .pipe(map(request => RequestDto.fromJson(request, this.generateColorService)))
            .pipe(finalize(() => observer.complete()))
            .subscribe(request => observer.next(request));
        }, err => observer.error(err));
    });
  }

  public createRequest(requestPayload: RequestPayload): Observable<RequestDto> {
    return this.apiService.post(`${this.apiRequestEndpoint}/request`, requestPayload.toJson())
      .pipe(map(response => _.get(response, 'body', null)))
      .pipe(filter(body => !!body))
      .pipe(map(body => RequestDto.fromJson(body)));
  }

  public approveRequest(idRequest: number): Observable<boolean> {
    return this.apiService.patch(`${this.apiRequestEndpoint}/request/approve/${idRequest}`)
      .pipe(map(response => _.get(response, 'body', null)))
      .pipe(filter(body => !!body))
      .pipe(map(body => RequestDto.fromJson(body)))
      .pipe(map(request => request.status === RequestStatus.APPROVE));
  }

  public rejectRequest(idRequest: number): Observable<boolean> {
    return this.apiService.patch(`${this.apiRequestEndpoint}/request/reject/${idRequest}`)
      .pipe(map(response => _.get(response, 'body', null)))
      .pipe(filter(body => !!body))
      .pipe(map(body => RequestDto.fromJson(body)))
      .pipe(map(request => request.status === RequestStatus.REJECT));
  }

  // endregion
}
