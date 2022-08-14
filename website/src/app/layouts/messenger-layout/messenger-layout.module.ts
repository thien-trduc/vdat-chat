import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MessengerLayoutComponent } from './messenger-layout.component';
import {MessengerLayoutRouting} from './messenger-layout.routing';
import {NzLayoutModule} from 'ng-zorro-antd/layout';
import {NzResizableModule} from 'ng-zorro-antd/resizable';
import {MessengerModule} from '../../pages/messenger/messenger.module';
import {NzGridModule} from 'ng-zorro-antd/grid';
import {NzAvatarModule} from 'ng-zorro-antd/avatar';
import {NzDropDownModule} from 'ng-zorro-antd/dropdown';
import {NzMenuModule} from 'ng-zorro-antd/menu';
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import {NzIconModule} from "ng-zorro-antd/icon";
import {TranslateModule} from "@ngx-translate/core";
import {NzModalModule} from "ng-zorro-antd/modal";
import {AppModule} from "../../app.module";
import {UserModule} from "../../pages/user/user.module";
import {NzButtonModule} from "ng-zorro-antd/button";
import {NzBadgeModule} from "ng-zorro-antd/badge";
import {NzSpinModule} from "ng-zorro-antd/spin";
import {NzResultModule} from "ng-zorro-antd/result";

@NgModule({
  declarations: [
    MessengerLayoutComponent
  ],
  imports: [
    CommonModule,
    MessengerLayoutRouting,
    NzLayoutModule,
    NzResizableModule,
    NzGridModule,
    NzAvatarModule,
    NzDropDownModule,
    MessengerModule,
    NzMenuModule,
    NzTypographyModule,
    NzIconModule,
    NzModalModule,
    UserModule,
    TranslateModule,
    NzButtonModule,
    NzBadgeModule,
    NzSpinModule,
    NzResultModule
  ]
})
export class MessengerLayoutModule { }
