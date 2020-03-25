import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceViewComponent } from './view/view.component';
import { ResourceCreateComponent } from './create/create.component';
import {ReactiveFormsModule} from '@angular/forms';



@NgModule({
  declarations: [ResourceViewComponent, ResourceCreateComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule
  ]
})
export class ResourceModule { }
