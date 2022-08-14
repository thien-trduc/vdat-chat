import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {InviteComponent} from '../../pages/invite/invite.component';

const routes: Routes = [
  {
    path: '',
    component: InviteComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class InviteLayoutRouting {
}
