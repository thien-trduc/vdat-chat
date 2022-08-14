import {Component, Input, OnChanges, OnInit, SimpleChanges} from '@angular/core';
import {Group} from '../../../core/models/group/group.model';
import {FileService} from '../../../core/services/collectors/file.service';
import {BehaviorSubject} from 'rxjs';
import {UploadFileDto} from '../../../core/models/file/upload-file.dto';
import * as _ from 'lodash';
import {filter} from 'rxjs/operators';
import {MessageType} from "../../../core/models/message/message-type.enum";
import {NzMessageService} from "ng-zorro-antd/message";

@Component({
  selector: 'app-message-files-shared-page',
  templateUrl: './message-files-shared-page.component.html',
  styleUrls: ['./message-files-shared-page.component.scss']
})
export class MessageFilesSharedPageComponent implements OnInit, OnChanges {

  @Input() currentGroup: Group;

  public listFile = new Array<UploadFileDto>();
  public loading = new BehaviorSubject<boolean>(false);

  constructor(private fileService: FileService,
              private nzMessageService: NzMessageService) {
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!changes.currentGroup && !!this.currentGroup) {
      this.listFile = new Array<UploadFileDto>();

      this.loading.next(true);
      this.fileService.getListFileOfGroup(this.currentGroup.id)
        .pipe(filter(file => !this.listFile.find(iter => iter.nameFile === file.nameFile)))
        .subscribe(file => {
            this.listFile.push(file);
            this.loading.next(false);
          }, err => this.loading.next(false),
          () => this.loading.next(false));
    }
  }

  public onDownloadFile(file: UploadFileDto): void {
    if (!!file) {
      const groupId = this.currentGroup.id;

      this.loading.next(true);
      this.fileService.getDownloadLink(groupId, file.nameFile)
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
