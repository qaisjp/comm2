import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceViewComponent } from './view/view.component';
import { ResourceCreateComponent } from './create/create.component';
import {ReactiveFormsModule} from '@angular/forms';
import {MomentModule} from 'ngx-moment';



@NgModule({
  declarations: [ResourceViewComponent, ResourceCreateComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MomentModule
  ]
})
export class ResourceModule { }
