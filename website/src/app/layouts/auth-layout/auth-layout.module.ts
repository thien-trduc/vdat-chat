import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthLayoutComponent } from './auth-layout.component';
import {AuthLayoutRouting} from './auth-layout.routing';
import {AuthenticateModule} from '../../pages/authenticate/authenticate.module';

@NgModule({
  declarations: [
    AuthLayoutComponent
  ],
  imports: [
    CommonModule,
    AuthenticateModule,
    AuthLayoutRouting
  ]
})
export class AuthLayoutModule { }
