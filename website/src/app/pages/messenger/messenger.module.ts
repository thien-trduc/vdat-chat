import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {MessageListPageComponent} from './message-list-page/message-list-page.component';
import {MessageContentPageComponent} from './message-content-page/message-content-page.component';
import {MessageInfoPageComponent} from './message-info-page/message-info-page.component';
import {NzGridModule} from 'ng-zorro-antd/grid';
import {NzTypographyModule} from 'ng-zorro-antd/typography';
import {NzButtonModule} from 'ng-zorro-antd/button';
import {NzIconModule} from 'ng-zorro-antd/icon';
import {NzToolTipModule} from 'ng-zorro-antd/tooltip';
import {NzInputModule} from 'ng-zorro-antd/input';
import {NzFormModule} from 'ng-zorro-antd/form';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {NzListModule} from 'ng-zorro-antd/list';
import {NzAvatarModule} from 'ng-zorro-antd/avatar';
import {RouterModule} from '@angular/router';
import {NzLayoutModule} from 'ng-zorro-antd/layout';
import {NzSpaceModule} from 'ng-zorro-antd/space';
import {NzDropDownModule} from 'ng-zorro-antd/dropdown';
import {NzCollapseModule} from 'ng-zorro-antd/collapse';
import {NzPipesModule} from 'ng-zorro-antd/pipes';
import {NzDividerModule} from 'ng-zorro-antd/divider';
import {MessageEmptyPageComponent} from './message-empty-page/message-empty-page.component';
import {NzResultModule} from 'ng-zorro-antd/result';
import {NzBadgeModule} from 'ng-zorro-antd/badge';
import {NzModalModule} from 'ng-zorro-antd/modal';
import {UserModule} from '../user/user.module';
import {MessageModule} from '../../components/message/message.module';
import {NzMessageModule, NzMessageService} from 'ng-zorro-antd/message';
import {MessageCreatePageComponent} from './message-create-page/message-create-page.component';
import {NzUploadModule} from 'ng-zorro-antd/upload';
import {NzTagModule} from 'ng-zorro-antd/tag';
import {NzSelectModule} from 'ng-zorro-antd/select';
import {NzCheckboxModule} from 'ng-zorro-antd/checkbox';
import {ScrollingModule} from '@angular/cdk/scrolling';
import {NzSkeletonModule} from 'ng-zorro-antd/skeleton';
import {CustomPipesModule} from '../../pipes/custom-pipes.module';
import {MessageListMemberPageComponent} from './message-list-member-page/message-list-member-page.component';
import {NzPopconfirmModule} from 'ng-zorro-antd/popconfirm';
import {MessageEditInfoPageComponent} from './message-edit-info-page/message-edit-info-page.component';
import {MessageImagesSharedPageComponent} from './message-images-shared-page/message-images-shared-page.component';
import {MessageFilesSharedPageComponent} from './message-files-shared-page/message-files-shared-page.component';
import {NzDrawerModule} from 'ng-zorro-antd/drawer';
import {NzImageModule} from 'ng-zorro-antd/image';
import {GroupRolePipe} from '../../pipes/group/group-role.pipe';
import {UploadModule} from '../../components/upload/upload.module';
import {MessageQrCodePageComponent} from './message-qr-code-page/message-qr-code-page.component';
import {UserInfoPipe} from '../../pipes/user-info.pipe';
import {MessageListPendingMembersPageComponent} from './message-list-pending-members-page/message-list-pending-members-page.component';
import {MessageReplyContentPageComponent} from './message-reply-content-page/message-reply-content-page.component';
import {NzNotificationModule} from 'ng-zorro-antd/notification';
import {NzPopoverModule} from 'ng-zorro-antd/popover';
import {PickerModule} from '@ctrl/ngx-emoji-mart';
import {NzAffixModule} from "ng-zorro-antd/affix";
import {NzBackTopModule} from "ng-zorro-antd/back-top";
import {NzSpinModule} from "ng-zorro-antd/spin";

@NgModule({
  declarations: [
    MessageListPageComponent,
    MessageContentPageComponent,
    MessageInfoPageComponent,
    MessageEmptyPageComponent,
    MessageCreatePageComponent,
    MessageListMemberPageComponent,
    MessageEditInfoPageComponent,
    MessageImagesSharedPageComponent,
    MessageFilesSharedPageComponent,
    MessageQrCodePageComponent,
    MessageListPendingMembersPageComponent,
    MessageReplyContentPageComponent,
  ],
  exports: [
    MessageListPageComponent,
    MessageContentPageComponent,
    MessageInfoPageComponent,
    MessageEmptyPageComponent
  ],
    imports: [
        CommonModule,
        NzGridModule,
        NzTypographyModule,
        NzButtonModule,
        NzIconModule,
        NzToolTipModule,
        NzInputModule,
        NzFormModule,
        ReactiveFormsModule,
        NzListModule,
        NzAvatarModule,
        RouterModule,
        NzLayoutModule,
        NzSpaceModule,
        NzDropDownModule,
        NzCollapseModule,
        NzPipesModule,
        NzDividerModule,
        NzResultModule,
        NzBadgeModule,
        NzModalModule,
        UserModule,
        MessageModule,
        ScrollingModule,
        NzMessageModule,
        NzUploadModule,
        NzTagModule,
        NzSelectModule,
        NzCheckboxModule,
        NzSkeletonModule,
        CustomPipesModule,
        FormsModule,
        NzPopconfirmModule,
        NzDrawerModule,
        NzImageModule,
        UploadModule,
        NzNotificationModule,
        NzPopoverModule,
        PickerModule,
        NzAffixModule,
        NzBackTopModule,
        NzSpinModule
    ],
  providers: [
    GroupRolePipe,
    UserInfoPipe
  ]
})
export class MessengerModule {
}
