import {NgModule} from '@angular/core';
import { CommonModule } from '@angular/common';
import { ResourceLayoutComponent } from './layout/layout.component';
import { ResourceCreateComponent } from './create/create.component';
import {ReactiveFormsModule} from '@angular/forms';
import {MomentModule} from 'ngx-moment';
import {RouterModule} from '@angular/router';
import { ResourceManageComponent } from './manage/manage.component';
import {OcticonModule} from '../octicon/octicon.module';
import {ResourceRoutingModule} from './resource-routing.module';
import {ResourceAboutComponent} from './about/about.component';
import {ResourceVersionsComponent} from './versions/versions.component';
import {ResourcePeopleComponent} from './people/people.component';
import { ResourceReviewsComponent } from './reviews/reviews.component';
import { ResourceUploadComponent } from './upload/upload.component';


@NgModule({
  declarations: [ResourceLayoutComponent, ResourceCreateComponent, ResourceManageComponent, ResourceAboutComponent, ResourceVersionsComponent, ResourcePeopleComponent, ResourceReviewsComponent, ResourceUploadComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MomentModule,
    RouterModule,
    OcticonModule,
    ResourceRoutingModule,
  ]
})
export class ResourceModule { }
