import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ProfileComponent } from './profile.component';
import {ResourceViewComponent} from '../resource/view/view.component';
import {ResourceManageComponent} from '../resource/manage/manage.component';


const routes: Routes = [
  {
    path: ':username',
    component: ProfileComponent,
  },
  {
    path: ':username/:resource',
    component: ResourceViewComponent,
  },
  {
    path: ':username/:resource/manage',
    component: ResourceManageComponent,
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ProfileRoutingModule { }
