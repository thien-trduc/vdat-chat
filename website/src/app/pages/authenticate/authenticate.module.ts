import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {AuthPageComponent} from './auth-page/auth-page.component';
import {NzGridModule} from 'ng-zorro-antd/grid';
import {NzTypographyModule} from 'ng-zorro-antd/typography';
import {NzAvatarModule} from 'ng-zorro-antd/avatar';
import {NzButtonModule} from 'ng-zorro-antd/button';
import {NzIconModule} from 'ng-zorro-antd/icon';
import {NzCardModule} from "ng-zorro-antd/card";

@NgModule({
  declarations: [
    AuthPageComponent
  ],
  imports: [
    CommonModule,
    NzGridModule,
    NzTypographyModule,
    NzAvatarModule,
    NzButtonModule,
    NzIconModule,
    NzCardModule
  ]
})
export class AuthenticateModule {
}
