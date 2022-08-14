import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {MessengerLayoutComponent} from './messenger-layout.component';

const routes: Routes = [
  {
    path: '',
    component: MessengerLayoutComponent
  },
  {
    path: ':id',
    component: MessengerLayoutComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MessengerLayoutRouting {
}
