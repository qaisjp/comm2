import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ProfileRoutingModule } from './profile-routing.module';
import { ProfileComponent } from './profile.component';
import {MomentModule} from 'ngx-moment';
import {OcticonModule} from '../octicon/octicon.module';
import { SettingsLayoutComponent } from '../settings/layout.component';


@NgModule({
  declarations: [ProfileComponent, SettingsLayoutComponent],
  imports: [
    CommonModule,
    ProfileRoutingModule,
    MomentModule,
    OcticonModule,
  ]
})
export class ProfileModule { }
