import {Injectable} from '@angular/core';
import {ApiService} from '../commons/api.service';
import {from, Observable} from 'rxjs';
import {environment} from '../../../../environments/environment';
import {NzUploadFile} from 'ng-zorro-antd/upload';
import {filter, map, switchMap} from 'rxjs/operators';
import * as _ from 'lodash';
import {UploadFileDto} from '../../models/file/upload-file.dto';
import {HttpClient} from '@angular/common/http';
import {FileType} from '../../constants/file-type.const';

@Injectable({
  providedIn: 'root'
})
export class FileService {
  private apiEndpoint = environment.api.files;

  constructor(private apiService: ApiService,
              private httpClient: HttpClient) {
  }

  public uploadFileInGroup(groupId: number, file: File | NzUploadFile): Observable<UploadFileDto> {
    const url = `${this.apiEndpoint}/upload/${groupId}`;

    return this.apiService.postFile(url, file)
      .pipe(filter(response => !!response))
      .pipe(map(response => _.get(response, 'body', {})))
      .pipe(map(body => UploadFileDto.fromJson(body)));
  }

  public getDownloadLink(groupId: number, fileName: string): Observable<string> {
    const url = `${this.apiEndpoint}/download/${groupId}/${encodeURIComponent(fileName)}`;

    return this.apiService.get(url)
      .pipe(filter(response => !!response))
      .pipe(map(response => _.get(response, 'body', {})))
      .pipe(map(body => UploadFileDto.fromJson(body)))
      .pipe(map(fileInfo => fileInfo.fileUrl));
  }

  public getListImageOfGroup(groupId: number): Observable<UploadFileDto> {
    return this.getFilesOfGroup(groupId, FileType.MEDIA)
      .pipe(filter(uploadFileDto => !!uploadFileDto));
  }

  public getListFileOfGroup(groupId: number): Observable<UploadFileDto> {
    return this.getFilesOfGroup(groupId, FileType.FILE)
      .pipe(filter(uploadFileDto => !!uploadFileDto));
  }

  private getFilesOfGroup(groupId: number, fileType: FileType): Observable<UploadFileDto> {
    const url = `${this.apiEndpoint}/${groupId}?type=${fileType}`;

    return this.apiService.get(url)
      .pipe(filter(response => !!response))
      .pipe(map(response => _.get(response, 'body', [])))
      .pipe(switchMap(arr => from(arr)))
      .pipe(map(body => UploadFileDto.fromJson(body)));
  }
}
