import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {AuthPageComponent} from '../../pages/authenticate/auth-page/auth-page.component';
import {AuthLayoutComponent} from './auth-layout.component';

const routes: Routes = [
  {
    path: '',
    component: AuthLayoutComponent,
    children: [
      {
        path: '',
        component: AuthPageComponent
      }
    ]
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AuthLayoutRouting {
}
