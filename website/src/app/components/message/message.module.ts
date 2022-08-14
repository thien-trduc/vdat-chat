import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MessageComponent} from './message.component';
import {MessageTextComponent} from './message-text/message-text.component';
import {MessageImageComponent} from './message-image/message-image.component';
import {MessageVideoComponent} from './message-video/message-video.component';
import {NzGridModule} from 'ng-zorro-antd/grid';
import {NzAvatarModule} from 'ng-zorro-antd/avatar';
import {NzSpaceModule} from 'ng-zorro-antd/space';
import {NzButtonModule} from 'ng-zorro-antd/button';
import {NzIconModule} from 'ng-zorro-antd/icon';
import {NzToolTipModule} from 'ng-zorro-antd/tooltip';
import {NzListModule} from 'ng-zorro-antd/list';
import {CustomPipesModule} from '../../pipes/custom-pipes.module';
import {NzBadgeModule} from 'ng-zorro-antd/badge';
import {NzDropDownModule} from 'ng-zorro-antd/dropdown';
import {NzTypographyModule} from 'ng-zorro-antd/typography';
import {NzImageModule} from 'ng-zorro-antd/image';
import {MessageInputComponent} from './message-input/message-input.component';
import {NzUploadModule} from 'ng-zorro-antd/upload';
import {NzPopoverModule} from 'ng-zorro-antd/popover';
import {PickerModule} from '@ctrl/ngx-emoji-mart';
import {ReactiveFormsModule} from '@angular/forms';
import {NzFormModule} from 'ng-zorro-antd/form';
import {NzInputModule} from 'ng-zorro-antd/input';
import {NzModalModule} from 'ng-zorro-antd/modal';
import {MessageFileComponent} from './message-file/message-file.component';
import {NzMessageModule} from 'ng-zorro-antd/message';
import {NzPipesModule} from 'ng-zorro-antd/pipes';


@NgModule({
  declarations: [
    MessageComponent,
    MessageTextComponent,
    MessageImageComponent,
    MessageVideoComponent,
    MessageInputComponent,
    MessageFileComponent
  ],
  exports: [
    MessageComponent,
    MessageInputComponent
  ],
  imports: [
    CommonModule,
    NzGridModule,
    NzAvatarModule,
    NzSpaceModule,
    NzButtonModule,
    NzIconModule,
    NzToolTipModule,
    NzListModule,
    CustomPipesModule,
    NzBadgeModule,
    NzDropDownModule,
    NzTypographyModule,
    NzImageModule,
    NzUploadModule,
    NzPopoverModule,
    PickerModule,
    ReactiveFormsModule,
    NzFormModule,
    NzInputModule,
    NzModalModule,
    NzMessageModule,
    NzPipesModule
  ]
})
export class MessageModule {
}
