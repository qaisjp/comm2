import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { LoginPageComponent } from './auth/login-page/login-page.component';
import {RegisterPageComponent} from './auth/register-page/register-page.component';
import {ResourceViewComponent} from './resource/view/view.component';
import {ResourceCreateComponent} from './resource/create/create.component';


const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
    pathMatch: 'full'
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
    loadChildren: 'src/app/admin/admin.module#AdminModule'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
