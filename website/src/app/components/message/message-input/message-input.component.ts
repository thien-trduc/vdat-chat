import {Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {User} from '../../../core/models/user.model';
import {NzUploadFile, NzUploadXHRArgs} from 'ng-zorro-antd/upload';
import {HttpEvent, HttpEventType, HttpResponse} from '@angular/common/http';
import {UploadFileDto} from '../../../core/models/file/upload-file.dto';
import {AbstractControl, FormControl, FormGroup, Validators} from '@angular/forms';
import {Message} from '../../../core/models/message/message.model';
import * as _ from 'lodash';
import {EmojiEvent} from '@ctrl/ngx-emoji-mart/ngx-emoji';
import {MessageService} from '../../../core/services/ws/message.service';
import {ApiService} from '../../../core/services/commons/api.service';
import {environment} from '../../../../environments/environment';
import {BehaviorSubject, Observable} from 'rxjs';
import {NzMessageService} from "ng-zorro-antd/message";

@Component({
  selector: 'app-message-input',
  templateUrl: './message-input.component.html',
  styleUrls: ['./message-input.component.scss']
})
export class MessageInputComponent implements OnInit, OnChanges {

  @Input() currentGroup: Group;
  @Input() currentUser: User;
  @Input() parentMessage: Message;

  public formMessage: FormGroup;
  public apiEndpoint: string;
  public previewImage: string | undefined = '';
  public previewVisible = new BehaviorSubject<boolean>(false);
  public visibleEmojiPicker = new BehaviorSubject<boolean>(false);

  constructor(private messageService: MessageService,
              private apiService: ApiService,
              private nzMessageService: NzMessageService) {
    this.formMessage = this.createFormMessage();
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!changes.currentGroup && !!this.currentGroup) {
      this.apiEndpoint = `${environment.api.files}/upload/${this.currentGroup?.id}`;
    }
  }

  private getBase64(file: File): Promise<string | ArrayBuffer | null> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = error => reject(error);
    });
  }

  public async handlePreviewImage(file: NzUploadFile): Promise<void> {
    if (!!file && !file.url && !file.preview) {
      file.preview = await this.getBase64(file.originFileObj);
    }
    this.previewImage = file.url || file.preview;
    this.previewVisible.next(true);
  }

  // region Event
  public onSendMessage(event?: any): void {
    if (this.formMessage.valid && !!this.currentGroup) {
      if (!!event) {
        event.preventDefault();
      }

      const groupId: number = this.currentGroup.id;

      if (this.formMessage.invalid && !!!groupId) {
        return;
      }

      const formValue = this.formMessage.getRawValue();
      const message = _.get(formValue, 'message', '');

      this.messageService.sendTextMessage(message, groupId,
        !!this.parentMessage ? this.parentMessage.id : null);

      // clear input
      this.formMessage.patchValue({message: ''});
    }
  }

  public onSelectEmoji(emojiEvent: EmojiEvent): void {
    const messageControl = this.formMessage.get('message');
    const message = messageControl.value;
    const emoji = emojiEvent.emoji.native;

    messageControl.setValue(`${message}${!!message ? ' ' : ''}${emoji}`);
  }

  // endregion

  // region Send File
  uploadFile = (item: NzUploadXHRArgs) => {
    const formData = new FormData();
    formData.append('file', item.file as any);

    const fileName = item.file.name;
    const loadingMessage = this.nzMessageService.loading(
      `Đang gửi tệp ${fileName}`,
      {
        nzDuration: 0
      }
    ).messageId;

    return this.apiService.postFile(this.apiEndpoint, formData)
      .subscribe((event: HttpEvent<{}>) => {
        if (event.type === HttpEventType.UploadProgress) {
          if (event.total > 0) {
            (event as any).percent = event.loaded / event.total * 100;
          }
          item.onProgress(event, item.file);
        } else if (event instanceof HttpResponse) {
          const fileUploaded = UploadFileDto.fromJson(event.body);
          item.onSuccess(fileUploaded.fileUrl, item.file, event);

          // send message
          if (!!this.currentGroup) {
            const groupId: number = this.currentGroup.id;

            if (!!!groupId) {
              return;
            }

            this.messageService.sendFileMessage(fileUploaded.nameFile, groupId,
              !!this.parentMessage ? this.parentMessage.id : null);

            this.nzMessageService.remove(loadingMessage);
          }
        }
      }, (err) => {
        item.onError(err, item.file);
        this.nzMessageService.remove(loadingMessage);
        this.nzMessageService.error(`Lỗi gửi tệp ${fileName}`);
      });
  }

  uploadImage = (item: NzUploadXHRArgs) => {
    const formData = new FormData();
    formData.append('file', item.file as any);

    const fileName = item.file.name;
    const loadingMessage = this.nzMessageService.loading(
      `Đang gửi hình ảnh ${fileName}`,
      {
        nzDuration: 0
      }
    ).messageId;

    return this.apiService.postFile(this.apiEndpoint, formData)
      .subscribe((event: HttpEvent<{}>) => {
        if (event.type === HttpEventType.UploadProgress) {
          if (event.total > 0) {
            (event as any).percent = event.loaded / event.total * 100;
          }
          item.onProgress(event, item.file);
        } else if (event instanceof HttpResponse) {
          const fileUploaded = UploadFileDto.fromJson(event.body);
          item.onSuccess(fileUploaded.fileUrl, item.file, event);

          // send message
          if (!!this.currentGroup) {
            const groupId: number = this.currentGroup.id;

            if (!!!groupId) {
              return;
            }

            this.messageService.sendImageMessage(fileUploaded.nameFile, groupId,
              !!this.parentMessage ? this.parentMessage.id : null);

            this.nzMessageService.remove(loadingMessage);
          }
        }
      }, (err) => {
        item.onError(err, item.file);
        this.nzMessageService.remove(loadingMessage);
        this.nzMessageService.error(`Lỗi gửi hình ảnh ${fileName}`);
      });
  }
  // endregion

  // region Form
  public get messageControl(): AbstractControl {
    return this.formMessage.get('message');
  }

  private createFormMessage(message?: string): FormGroup {
    return new FormGroup({
      message: new FormControl(!!message ? message : '', Validators.required)
    });
  }

  // endregion
}
