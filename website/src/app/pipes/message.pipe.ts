import {Pipe, PipeTransform} from '@angular/core';
import {Message} from '../core/models/message/message.model';
import {User} from '../core/models/user.model';
import {MessageType} from '../core/models/message/message-type.enum';
import * as _ from 'lodash';
import {from, Observable, of} from 'rxjs';
import {filter, map, takeUntil} from 'rxjs/operators';
import {FileService} from '../core/services/collectors/file.service';
import {Group} from '../core/models/group/group.model';

@Pipe({
  name: 'message'
})
export class MessagePipe implements PipeTransform {
  constructor(private fileService: FileService) {
  }

  transform(message: Message, key: 'isSender' | 'lastMessage' | 'fileName' | 'iconFile' | 'imageUrl',
            currentUser?: User, currentGroup?: Group): any | Observable<string> {
    if (!!message && !!currentUser && key === 'isSender') {
      return message.sender.userId === currentUser.userId;
    } else {
      switch (key) {
        case 'lastMessage':
          return this.getLastMessage(message);
        case 'fileName':
          return this.getFileNameInMessage(message);
        case 'iconFile':
          return this.getIconFile(message);
        case 'imageUrl':
          return this.getImageUrlInMessage(message, currentGroup);
      }
    }
  }

  private getLastMessage(message: Message): string {
    if (!!message.deletedAt) {
      return 'tin nhắn đã bị thu hồi';
    } else {
      switch (message.messageType) {
        case MessageType.TEXT_MESSAGE:
          return message.message;
        case MessageType.FILE_MESSAGE:
          return 'đã gửi 1 tệp';
        case MessageType.IMAGE_MESSAGE:
          return 'đã gửi 1 hình ảnh';
      }
    }
  }

  private getImageUrlInMessage(message: Message, currentGroup: Group): Observable<string> {
    if (message.message.startsWith('http') || message.message.startsWith('https')) {
      return of(message.message);
    }

    return this.fileService.getDownloadLink(currentGroup.id, message.message)
      .pipe(map(fileInfo => fileInfo));
  }

  private getFileNameInMessage(message: Message): string {
    if (message.messageType === MessageType.FILE_MESSAGE) {
      const fileName: string = _.get(message, 'message', '').trim();
      const startIndex = Math.max(fileName.indexOf('_', 0), 0);
      return fileName.substring(startIndex > 0 ? startIndex + 1 : startIndex);
    }
    return message.message;
  }

  private getIconFile(message: Message): Observable<string> {
    const fileName = this.getFileNameInMessage(message);
    const extFile = fileName.split('.').pop().trim();

    return from(this.getListFileIcon())
      .pipe(filter(obj => !!obj))
      .pipe(filter(obj => !!obj.extensions.split('|').find(ext => ext === extFile)))
      .pipe(map(obj => obj.icon));
  }

  private getListFileIcon(): Array<{extensions: string, icon: string}> {
    return [
      {
        extensions: 'txt|tex',
        icon: 'file-text'
      },
      {
        extensions: 'gif',
        icon: 'file-gif'
      },
      {
        extensions: 'md',
        icon: 'file-markdown'
      },
      {
        extensions: 'pdf',
        icon: 'file-pdf'
      },
      {
        extensions: 'ai|bmp|ico|jpeg|jpg|png|ps|psd|svg|tif|tiff',
        icon: 'file-image'
      },
      {
        extensions: 'doc|docx|odt|rtf|wpd',
        icon: 'file-word'
      },
      {
        extensions: 'ods|xls|xlsm|xlsx',
        icon: 'file-excel'
      },
      {
        extensions: 'key|odp|ppt|pptx',
        icon: 'file-ppt'
      },
      {
        extensions: '7z|arj|rar|gz|z|zip',
        icon: 'file-zip'
      },
      {
        extensions: '*',
        icon: 'file'
      }
    ];
  }
}
