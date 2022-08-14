import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {filter, map} from 'rxjs/operators';
import {Invite} from '../../core/models/invite.model';
import {EncryptService} from '../../core/services/commons/encrypt.service';
import {NzMessageService} from 'ng-zorro-antd/message';
import {Group} from "../../core/models/group/group.model";
import {User} from "../../core/models/user.model";
import {StorageService} from "../../core/services/commons/storage.service";
import {GroupService} from "../../core/services/collectors/group.service";
import {RequestPayload} from "../../core/models/request/request.payload";
import {RequestDto} from "../../core/models/request/request.dto";
import {RequestStatus} from "../../core/constants/request-status.enum";

@Component({
  selector: 'app-invite',
  templateUrl: './invite.component.html',
  styleUrls: ['./invite.component.scss']
})
export class InviteComponent implements OnInit {

  public invite: Invite;
  public group: Group;
  public userInvite: User;
  public currentUser: User;
  public requestDto: RequestDto;

  constructor(private route: ActivatedRoute,
              private encryptService: EncryptService,
              private messageService: NzMessageService,
              private storageService: StorageService,
              private groupService: GroupService,
              private router: Router) {
    this.currentUser = this.storageService.userInfo;

    this.route.queryParams
      .pipe(filter(queries => !!queries))
      .pipe(filter(queries => !!queries.g))
      .pipe(map(queries => Invite.fromData(queries.g, this.encryptService)))
      .pipe(filter(invite => !!invite))
      .subscribe(invite => this.invite = invite,
        err => this.messageService.error(err),
        () => {
          if (!this.invite) {
            return;
          }
        });
  }

  ngOnInit(): void {
  }

  public onConfirm(): void {
    const requestPayload = new RequestPayload(
      this.invite.groupId,
      this.invite.userInviteId,
      this.currentUser.userId
    );

    this.groupService.createRequest(requestPayload)
      .subscribe(requestDto => {
        if (!!requestDto && requestDto.status === RequestStatus.APPROVE) {
          this.messageService.success('Tham gia nhóm thành công');
          this.router.navigate(['/messages', this.invite.groupId])
            .then(() => {
            });
        } else {
          this.messageService.error('Lỗi tham gia nhóm! Vui lòng thử lại sau');
        }
      });
  }
}
