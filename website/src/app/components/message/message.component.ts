import {Component, Input, OnInit, Output, EventEmitter} from '@angular/core';
import {Message} from '../../core/models/message/message.model';
import {User} from '../../core/models/user.model';
import {MessageService} from '../../core/services/ws/message.service';
import {NzModalService} from 'ng-zorro-antd/modal';
import {MessageType} from '../../core/models/message/message-type.enum';
import {Group} from '../../core/models/group/group.model';
import {BehaviorSubject} from 'rxjs';
import {FileService} from '../../core/services/collectors/file.service';
import {NzMessageService} from 'ng-zorro-antd/message';

@Component({
  selector: 'app-message',
  templateUrl: './message.component.html',
  styleUrls: ['./message.component.scss']
})
export class MessageComponent implements OnInit {

  @Input() message: Message;
  @Input() currentUser: User;
  @Input() currentGroup: Group;
  @Input() isViewInChildMessage: boolean;

  @Output() replyMessageEvent = new EventEmitter();

  public messageType = MessageType;
  public loading = new BehaviorSubject<boolean>(false);

  constructor(private messageService: MessageService,
              private modalService: NzModalService,
              private fileService: FileService,
              private nzMessageService: NzMessageService) {
  }

  ngOnInit(): void {
  }

  public onReplyMessage(): void {
    this.replyMessageEvent.emit();
  }

  public onDeleteMessage(): void {
    if (this.message.sender.userId === this.currentUser.userId) {
      this.modalService.confirm({
        nzTitle: 'Cảnh báo',
        nzContent: 'Bạn có chắc muốn xoá tin nhắn này không ?',
        nzOkText: 'Xoá tin nhắn',
        nzOkType: 'danger',
        nzCancelText: 'Huỷ',
        nzCentered: true,
        nzOnOk: () => this.messageService.deleteMessage(this.message.id, this.message.groupId)
      });
    }
  }

  public onDownloadFile(): void {
    if (!!this.message && this.message.messageType === MessageType.FILE_MESSAGE) {
      const fileName = this.message.message;
      const groupId = this.message.groupId;

      this.loading.next(true);
      this.fileService.getDownloadLink(groupId, fileName)
        .subscribe(fileUrl => {
            window.open(fileUrl, '_blank');
          }, err => {
            this.loading.next(false);
            this.nzMessageService.error('Lỗi tải tệp tin !');
          },
          () => this.loading.next(false));
    }
  }
}
