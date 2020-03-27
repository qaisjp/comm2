import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ProfileComponent } from './profile.component';
import {ResourceViewComponent} from '../resource/view/view.component';


const routes: Routes = [
  {
    path: 'u/:username',
    component: ProfileComponent,
  },
  {
    path: 'u/:username/:resource',
    component: ResourceViewComponent,
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ProfileRoutingModule { }
