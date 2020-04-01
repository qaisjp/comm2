import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { AdminRoutingModule } from './admin-routing.module';
import { AdminDashboardComponent } from './dashboard/dashboard.component';
import { AdminLayoutComponent } from './layout/layout.component';
import { AdminUsersComponent } from './users/users.component';
import { AdminBansComponent } from './bans/bans.component';


@NgModule({
  declarations: [AdminDashboardComponent, AdminLayoutComponent, AdminUsersComponent, AdminBansComponent],
  imports: [
    CommonModule,
    AdminRoutingModule
  ]
})
export class AdminModule { }
