import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ProfileComponent } from './profile.component';
import {ResourceLayoutComponent} from '../resource/layout/layout.component';
import {ResourceManageComponent} from '../resource/manage/manage.component';


const routes: Routes = [
  {
    path: ':username',
    component: ProfileComponent,
  },
  {
    path: ':username/:resource',
    loadChildren: 'src/app/resource/resource.module#ResourceModule'
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ProfileRoutingModule { }
