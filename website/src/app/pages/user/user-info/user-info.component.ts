import {Component, OnInit, EventEmitter, Input, Output, OnChanges, SimpleChanges} from '@angular/core';
import {User} from '../../../core/models/user.model';
import {FormControl, FormGroup} from '@angular/forms';
import {NzUploadXHRArgs} from "ng-zorro-antd/upload";
import {HttpEvent, HttpEventType, HttpResponse} from "@angular/common/http";
import {UploadFileDto} from "../../../core/models/file/upload-file.dto";
import {ApiService} from "../../../core/services/commons/api.service";
import {environment} from "../../../../environments/environment";
import {FileService} from "../../../core/services/collectors/file.service";

@Component({
  selector: 'app-user-info',
  templateUrl: './user-info.component.html',
  styleUrls: ['./user-info.component.scss']
})
export class UserInfoComponent implements OnInit, OnChanges {

  @Input() currentUser: User;
  @Output() currentUserChange: EventEmitter<User> = new EventEmitter<User>();

  public formUserInfo: FormGroup;
  public apiEndpoint = `${environment.api.files}/avatar/user`;

  constructor(private apiService: ApiService,
              private fileService: FileService) {
    this.formUserInfo = this.createFormGroup();
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.currentUser && !!this.currentUser) {
      this.formUserInfo = this.createFormGroup(this.currentUser);
    }
  }

  ngOnInit(): void {
  }

  uploadAvatar = (item: NzUploadXHRArgs) => {
    const formData = new FormData();
    formData.append('file', item.file as any);

    return this.apiService.patchFile(this.apiEndpoint, formData)
      .subscribe((event: HttpEvent<{}>) => {
        if (event.type === HttpEventType.UploadProgress) {
          if (event.total > 0) {
            (event as any).percent = event.loaded / event.total * 100;
          }
          item.onProgress(event, item.file);
        } else if (event instanceof HttpResponse) {
          const fileUploaded = UploadFileDto.fromJson(event.body);
          item.onSuccess(fileUploaded.fileUrl, item.file, event);
          this.currentUser.avatar = fileUploaded.fileUrl;
        }
      }, (err) => {
        item.onError(err, item.file);
      });
  }

  private createFormGroup(user?: User): FormGroup {
    const formGroup: FormGroup = new FormGroup({
      firstName: new FormControl(''),
      lastName: new FormControl(''),
      username: new FormControl({value: '', disabled: true})
    });

    if (!!user) {
      formGroup.patchValue(user);
    }

    return formGroup;
  }
}
