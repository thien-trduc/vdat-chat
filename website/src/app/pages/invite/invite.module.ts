import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {InviteComponent} from './invite.component';
import {NzResultModule} from 'ng-zorro-antd/result';
import {NzTypographyModule} from 'ng-zorro-antd/typography';
import {NzIconModule} from 'ng-zorro-antd/icon';
import {NzButtonModule} from 'ng-zorro-antd/button';
import {CustomPipesModule} from "../../pipes/custom-pipes.module";
import {NzAvatarModule} from "ng-zorro-antd/avatar";
import {NzGridModule} from "ng-zorro-antd/grid";
import {RouterModule} from "@angular/router";

@NgModule({
  declarations: [
    InviteComponent
  ],
  imports: [
    CommonModule,
    NzResultModule,
    NzTypographyModule,
    NzIconModule,
    NzButtonModule,
    CustomPipesModule,
    NzAvatarModule,
    NzGridModule,
    RouterModule
  ]
})
export class InviteModule {
}
