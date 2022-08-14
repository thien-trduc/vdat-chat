import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {AuthGuard} from './guards/auth.guard';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: 'messages',
  },
  {
    path: 'messages',
    loadChildren: () => import('./layouts/messenger-layout/messenger-layout.module').then(m => m.MessengerLayoutModule),
    canActivate: [AuthGuard],
  },
  {
    path: 'auth',
    loadChildren: () => import('./layouts/auth-layout/auth-layout.module').then(m => m.AuthLayoutModule),
  },
  {
    path: 'invite',
    loadChildren: () => import('./layouts/invite-layout/invite-layout.module').then(m => m.InviteLayoutModule),
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, {
    initialNavigation: 'enabledNonBlocking',
    onSameUrlNavigation: 'ignore',
    relativeLinkResolution: 'legacy',
    useHash: false
  })],
  exports: [RouterModule]
})
export class AppRouting {
}
