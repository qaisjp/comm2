import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {AdminDashboardComponent} from './dashboard/dashboard.component';
import {AdminLayoutComponent} from './layout/layout.component';
import {AdminUsersComponent} from './users/users.component';
import {AdminBansComponent} from './bans/bans.component';


const routes: Routes = [
  {
    path: 'dashboard',
    component: AdminDashboardComponent,
    data: {title: 'Dashboard'}
  },
  {
    path: 'users',
    component: AdminUsersComponent,
    data: {title: 'Users'},
  },
  {
    path: 'bans',
    component: AdminBansComponent,
    data: {title: 'Bans'},
  }
];

@NgModule({
  imports: [
    RouterModule.forChild([
      {
        path: '',
        component: AdminLayoutComponent,
        children: routes,
      }
    ]),
  ],
  exports: [RouterModule]
})
export class AdminRoutingModule { }
