import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceViewComponent } from './view/view.component';
import { ResourceCreateComponent } from './create/create.component';
import {ReactiveFormsModule} from '@angular/forms';
import {MomentModule} from 'ngx-moment';
import {RouterModule} from '@angular/router';



@NgModule({
  declarations: [ResourceViewComponent, ResourceCreateComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MomentModule,
    RouterModule
  ]
})
export class ResourceModule { }
