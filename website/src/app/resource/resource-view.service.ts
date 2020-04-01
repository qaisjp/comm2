import { Injectable } from '@angular/core';
import {Observable, of, ReplaySubject} from 'rxjs';
import {Resource, ResourcePackage, ResourceService} from './resource.service';
import {map, single, tap} from 'rxjs/operators';

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

  getKeyCounter(key: string): Observable<number> {
    if (key === 'people') {
      return this.resource$.pipe(map(r => r.authors.length));
    } else if (key === 'reviews') {

    } else if (key === 'versions') {
      return this.packages$.pipe(map(ps => ps.length));
    }
    return of(1337);
  }

}
