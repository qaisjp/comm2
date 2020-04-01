import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceViewComponent } from './view/view.component';
import { ResourceCreateComponent } from './create/create.component';
import {ReactiveFormsModule} from '@angular/forms';
import {MomentModule} from 'ngx-moment';
import {RouterModule} from '@angular/router';
import { ResourceManageComponent } from './manage/manage.component';
import {OcticonModule} from '../octicon/octicon.module';



@NgModule({
  declarations: [ResourceViewComponent, ResourceCreateComponent, ResourceManageComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MomentModule,
    RouterModule,
    OcticonModule
  ]
})
export class ResourceModule { }
