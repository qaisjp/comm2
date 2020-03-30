import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {ProfileComponent} from '../profile/profile.component';
import {AdminDashboardComponent} from './dashboard/dashboard.component';


const routes: Routes = [
  {
    path: 'dashboard',
    component: AdminDashboardComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AdminRoutingModule { }
