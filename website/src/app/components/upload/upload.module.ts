import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {NzUploadModule} from 'ng-zorro-antd/upload';
import {NzIconModule} from 'ng-zorro-antd/icon';
import {UploadImageComponent} from './upload-image/upload-image.component';
import {NzMessageModule} from 'ng-zorro-antd/message';

@NgModule({
  declarations: [
    UploadImageComponent
  ],
  exports: [
    UploadImageComponent
  ],
  imports: [
    CommonModule,
    NzUploadModule,
    NzIconModule,
    NzMessageModule
  ]
})
export class UploadModule {
}
