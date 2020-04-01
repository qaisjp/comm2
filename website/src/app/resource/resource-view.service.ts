import { Injectable } from '@angular/core';
import {Observable, ReplaySubject} from 'rxjs';
import {Resource, ResourcePackage, ResourceService} from './resource.service';
import {single, tap} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ResourceViewService {

  public resource$: ReplaySubject<Resource> = new ReplaySubject(1);
  public packages$: ReplaySubject<ResourcePackage[]> = new ReplaySubject(1);
  public downloadable = false;

  constructor(
    private resources: ResourceService,
  ) {
    this.resource$.subscribe((data: Resource) => {
      this.resources.getPackages(data.author_id, data.id).pipe(
        single(),
        tap(packages => {
          this.downloadable = packages.some(pkg => !pkg.draft);
        }),
      ).subscribe(this.packages$);
    });
  }

}
