import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {InviteLayoutComponent} from './invite-layout.component';
import {InviteLayoutRouting} from './invite-layout.routing';
import {InviteModule} from '../../pages/invite/invite.module';

@NgModule({
  declarations: [
    InviteLayoutComponent
  ],
  imports: [
    CommonModule,
    InviteModule,
    InviteLayoutRouting
  ]
})
export class InviteLayoutModule {
}
