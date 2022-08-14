import {Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges} from '@angular/core';
import {environment} from '../../../../environments/environment';
import {UploadFileDto} from '../../../core/models/file/upload-file.dto';
import {NzUploadFile, NzUploadXHRArgs} from 'ng-zorro-antd/upload';
import {HttpEvent, HttpEventType, HttpResponse} from '@angular/common/http';
import {ApiService} from '../../../core/services/commons/api.service';
import {FileService} from '../../../core/services/collectors/file.service';
import {BehaviorSubject, Observable, Observer} from 'rxjs';
import {NzMessageService} from 'ng-zorro-antd/message';
import {NgxImageCompressService} from 'ngx-image-compress';

@Component({
  selector: 'app-upload-image',
  templateUrl: './upload-image.component.html',
  styleUrls: ['./upload-image.component.scss']
})
export class UploadImageComponent implements OnInit, OnChanges {

  @Input() apiEndpoint = `${environment.api.files}/upload`;
  @Input() oldImage: string;
  @Output() fileName = new EventEmitter<string>();

  public thumbnailUrl: string;
  public loading = new BehaviorSubject<boolean>(false);

  file: any;
  predictions: number[];
  imageDataEvent: any;
  localUrl: any;
  localCompressedURl: any;
  isFlip = false;
  sizeOfOriginalImage: number;
  sizeOFCompressedImage: number;
  imgResultBeforeCompress: string;
  imgResultAfterCompress: string;

  constructor(private apiService: ApiService,
              private fileService: FileService,
              private nzMessageService: NzMessageService,
              private imageCompress: NgxImageCompressService) {
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges): void {
    if (!!changes && !!this.oldImage) {
      this.thumbnailUrl = this.oldImage;
    }
  }

  uploadImage = (item: NzUploadXHRArgs) => {
    const formData = new FormData();
    formData.append('file', item.file as any);

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
          this.fileName.emit(fileUploaded.nameFile);
        }
      }, (err) => {
        item.onError(err, item.file);
      });
  }

  transformFile = (file: NzUploadFile) => {
    return new Observable((observer: Observer<Blob>) => {
      this.getBase64(file as any, (imageBase64Str: string) => {
        this.imageCompress.compressFile(imageBase64Str, -1, 50, 50)
          .then(result => {
              const fileCompress = this.convertDataURItoBlob(result);
              console.log(fileCompress);
              observer.next(fileCompress);
              observer.complete();
            }
          );
      });
    });
  }

  public handleChange(info: { file: NzUploadFile }): void {
    switch (info.file.status) {
      case 'uploading':
        this.loading.next(true);
        break;
      case 'done':
        // tslint:disable-next-line:no-non-null-assertion
        this.getBase64(info.file!.originFileObj!, (img: string) => {
          this.loading.next(false);
          this.thumbnailUrl = img;
        });
        break;
      case 'error':
        this.nzMessageService.error('Lỗi tải lên hình ảnh');
        this.loading.next(false);
        break;
    }
  }

  private getBase64(img: File, callback: (img: string) => void): void {
    const reader = new FileReader();
    reader.addEventListener('load', () => callback(reader.result.toString()));
    reader.readAsDataURL(img);
  }

  private convertDataURItoBlob(dataURI): Blob {
    const byteString = atob(dataURI.split(',')[1]);

    const arrayBuffer = new ArrayBuffer(byteString.length);
    const int8Array = new Uint8Array(arrayBuffer);

    for (let i = 0; i < byteString.length; i++) {
      int8Array[i] = byteString.charCodeAt(i);
    }

    return new Blob([arrayBuffer], {type: 'image/jpeg'});
  }
}
