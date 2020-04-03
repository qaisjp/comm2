import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import {ResourceLayoutComponent} from './layout/layout.component';
import {ResourceManageComponent} from './manage/manage.component';
import {ResourceAboutComponent} from './about/about.component';
import {ResourceVersionsComponent} from './versions/versions.component';
import {ResourcePeopleComponent} from './people/people.component';
import {ResourceReviewsComponent} from './reviews/reviews.component';
import {ResourceUploadComponent} from './upload/upload.component';
import {ResourcePackageResolverService} from './package-resolver.service';


const routes: Routes = [
  {
    path: '',
    component: ResourceAboutComponent,
    data: {
      key: 'about',
      title: 'About',
      icon: 'repo',
    },
  },
  {
    path: 'versions',
    data: {
      key: 'versions',
      title: 'Versions',
      icon: 'bug',
      counter: true,
    },
    children: [
      {
        path: '',
        component: ResourceVersionsComponent,
      },
      {
        path: 'upload',
        component: ResourceUploadComponent,
      },
      {
        path: 'edit/:pkg_id',
        resolve: {pkg: ResourcePackageResolverService},
        component: ResourceUploadComponent,
      }
    ]
  },
  {
    path: 'reviews',
    component: ResourceReviewsComponent,
    data: {
      key: 'reviews',
      title: 'Reviews',
      icon: 'comment-discussion',
      counter: true,
    },
  },
  {
    path: 'people',
    component: ResourcePeopleComponent,
    data: {
      key: 'people',
      title: 'People',
      icon: 'person',
      counter: true,
    },
  },
  {
    path: 'manage',
    component: ResourceManageComponent,
    data: {
      key: 'manage',
      title: 'Settings',
      icon: 'gear',
      elevated: true,
    },
  },
];

@NgModule({
  imports: [
    RouterModule.forChild([
      {
        path: '',
        component: ResourceLayoutComponent,
        children: routes,
      }
    ]),
  ],
  exports: [RouterModule]
})
export class ResourceRoutingModule { }
