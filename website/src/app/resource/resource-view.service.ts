import { Injectable } from '@angular/core';
import {Observable, of, ReplaySubject, throwError} from 'rxjs';
import {Resource, ResourcePackage, ResourceService} from './resource.service';
import {catchError, last, map, single, switchMap, tap} from 'rxjs/operators';
import {HttpEvent, HttpEventType} from '@angular/common/http';
import {AlertService} from '../alert.service';
import {LogService} from '../log.service';

@Injectable({
  providedIn: 'root'
})
export class ResourceViewService {

  public resource$: ReplaySubject<Resource> = new ReplaySubject(1);
  public packages$: ReplaySubject<ResourcePackage[]> = new ReplaySubject(1);
  public downloadable = false;

  constructor(
    private resources: ResourceService,
    private alerts: AlertService,
    private log: LogService,
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

  private getDownloadEventMessage(event: HttpEvent<any>, pkg: ResourcePackage) {
    switch (event.type) {
      case HttpEventType.Sent:
        return `Downloading version ${pkg.version}.`;

      case HttpEventType.UploadProgress:
        // Compute and show the % done:
        const percentDone = Math.round(100 * event.loaded / event.total);
        return `Version "${pkg.version}" is ${percentDone}% downloaded.`;

      case HttpEventType.Response:
        return `Version "${pkg.version}" was completely download!`;

      default:
        return `Version "${pkg.version}" surprising download event: ${event.type}.`;
    }
  }

  download(pkg: ResourcePackage) {
    return this.resource$.pipe(
      switchMap(r => this.resources.download(r.author_id, r.id, pkg.id)),
      map(event => this.getDownloadEventMessage(event, pkg)),
      tap(message => this.alerts.setAlert(message)),
      last(),
      catchError(err => {
        this.alerts.setAlert(`Failed to download v${pkg.version}`);
        this.log.error(`failed to download v${pkg.version}`, err);
        return throwError(err);
      }),
    );
  }

}
