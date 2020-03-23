import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceViewComponent } from './view/view.component';
import { ResourceCreateComponent } from './create/create.component';



@NgModule({
  declarations: [ResourceViewComponent, ResourceCreateComponent],
  imports: [
    CommonModule
  ]
})
export class ResourceModule { }
