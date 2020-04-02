import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { LoginPageComponent } from './auth/login-page/login-page.component';
import {RegisterPageComponent} from './auth/register-page/register-page.component';
import {ResourceLayoutComponent} from './resource/layout/layout.component';
import {ResourceCreateComponent} from './resource/create/create.component';
import {AdminDashboardComponent} from './admin/dashboard/dashboard.component';
import {AdminLayoutComponent} from './admin/layout/layout.component';
import {SettingsLayoutComponent} from './settings/layout.component';
import {ProfileComponent} from './profile/profile.component';
import {SettingsProfileComponent} from './settings/profile.component';
import {SettingsAccountComponent} from './settings/account.component';


const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
    pathMatch: 'full',
    data: {navWide: true},
  },
  {
    path: 'login',
    component: LoginPageComponent,
    pathMatch: 'full',
  },
  {
    path: 'register',
    component: RegisterPageComponent,
    pathMatch: 'full',
  },
  {
    path: 'create',
    component: ResourceCreateComponent,
  },
  {
    path: 'u',
    loadChildren: 'src/app/profile/profile.module#ProfileModule'
  },
  {
    path: 'admin',
    loadChildren: 'src/app/admin/admin.module#AdminModule',
  },
  {
    path: 'settings',
    component: SettingsLayoutComponent,
    children: [
      {
        path: '',
        redirectTo: 'profile',
        pathMatch: 'full',
        data: {},
      },
      {
        path: 'profile',
        component: SettingsProfileComponent,
        data: {title: 'Profile'},
      },
      {
        path: 'account',
        component: SettingsAccountComponent,
        data: {title: 'Account'},
      },
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
