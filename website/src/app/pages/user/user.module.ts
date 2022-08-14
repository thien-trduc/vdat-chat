import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {UserInfoComponent} from './user-info/user-info.component';
import {NzGridModule} from "ng-zorro-antd/grid";
import {NzAvatarModule} from "ng-zorro-antd/avatar";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {NzButtonModule} from "ng-zorro-antd/button";
import {NzFormModule} from "ng-zorro-antd/form";
import {NzInputModule} from "ng-zorro-antd/input";
import { SearchUsersComponent } from './search-users/search-users.component';
import {NzIconModule} from "ng-zorro-antd/icon";
import {ScrollingModule} from "@angular/cdk/scrolling";
import {NzListModule} from "ng-zorro-antd/list";
import {NzSkeletonModule} from "ng-zorro-antd/skeleton";
import {NzTypographyModule} from "ng-zorro-antd/typography";
import {NzCheckboxModule} from "ng-zorro-antd/checkbox";
import {NzRadioModule} from "ng-zorro-antd/radio";
import {NzTagModule} from "ng-zorro-antd/tag";
import {AppModule} from "../../app.module";
import {CustomPipesModule} from "../../pipes/custom-pipes.module";
import {CustomDirectiveModule} from "../../directives/custom-directive.module";
import {NzBadgeModule} from "ng-zorro-antd/badge";
import {NzToolTipModule} from "ng-zorro-antd/tooltip";
import {NzUploadModule} from "ng-zorro-antd/upload";
import {NzSpinModule} from "ng-zorro-antd/spin";

@NgModule({
  declarations: [
    UserInfoComponent,
    SearchUsersComponent
  ],
  exports: [
    UserInfoComponent,
    SearchUsersComponent
  ],
    imports: [
        CommonModule,
        NzGridModule,
        NzAvatarModule,
        ReactiveFormsModule,
        NzButtonModule,
        NzFormModule,
        NzInputModule,
        FormsModule,
        NzIconModule,
        ScrollingModule,
        NzListModule,
        NzSkeletonModule,
        NzTypographyModule,
        NzCheckboxModule,
        NzRadioModule,
        NzTagModule,
        CustomPipesModule,
        CustomDirectiveModule,
        NzBadgeModule,
        NzToolTipModule,
        NzUploadModule,
        NzSpinModule
    ]
})
export class UserModule {
}
