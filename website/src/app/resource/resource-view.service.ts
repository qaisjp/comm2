import {Injectable} from '@angular/core';
import {EMPTY, Observable, of, ReplaySubject, throwError} from 'rxjs';
import {Resource, ResourcePackage, ResourceService} from './resource.service';
import {catchError, first, last, map, single, switchMap, takeUntil, takeWhile, tap} from 'rxjs/operators';
import {HttpEvent, HttpEventType} from '@angular/common/http';
import {AlertService} from '../alert.service';
import {LogService} from '../log.service';

@Injectable({
  providedIn: 'root'
})
export class ResourceViewService {

  constructor(
    private resources: ResourceService,
    private alerts: AlertService,
    private log: LogService,
  ) {
    this.resource$.subscribe((data: Resource) => {
      this.resources.getPackages(data.author_id, data.id).pipe(
        single(),
        tap(packages => {
          this.downloadable = packages.some(pkg => pkg.published_at);
        }),
      ).subscribe(this.packages$);
    });
  }

  public resource$: ReplaySubject<Resource> = new ReplaySubject(1);
  public packages$: ReplaySubject<ResourcePackage[]> = new ReplaySubject(1);
  public downloadable = false;

  downloadProgress: { [key: number]: number } = {}

  getKeyCounter(key: string): Observable<number> {
    if (key === 'people') {
      return this.resource$.pipe(map(r => r.authors.length));
    } else if (key === 'reviews') {

    } else if (key === 'versions') {
      return this.packages$.pipe(map(ps => ps.length));
    }
    return of(1337);
  }

  private getDownloadEventMessage(event: HttpEvent<any>, pkg: ResourcePackage): [boolean, string, Blob?] {
    switch (event.type) {
      case HttpEventType.Sent:
        this.downloadProgress[pkg.id] = 0;
        return [false, `Downloading version ${pkg.version}.`, null];

      case HttpEventType.ResponseHeader:
        return [false, `Receiving ${pkg.version}...`, null];

      case HttpEventType.DownloadProgress:
        // Compute and show the % done:
        const percentDone = Math.round(100 * event.loaded / event.total);
        this.downloadProgress[pkg.id] = percentDone;
        return [false, `Version "${pkg.version}" is ${percentDone}% downloaded.`, null];

      case HttpEventType.Response:
        delete this.downloadProgress[pkg.id];
        return [true, `Version "${pkg.version}" was completely downloaded!`, event.body];

      default:
        delete this.downloadProgress[pkg.id];
        return [true, `Version "${pkg.version}" surprising download event: ${event.type}.`, null];
    }
  }

  download(pkg: ResourcePackage): Observable<Blob> {
    return this.resource$.pipe(
      switchMap(r => this.resources.download(r.author_id, r.id, pkg.id)),
      map(event => this.getDownloadEventMessage(event, pkg)),
      tap(message => console.log(message[0], message[1], message[2])),
      // todo: last() was not working, so we have first() and this 'done' message[0] bool
      first(msg => msg[0] ),
      map(data => data[2]),
      catchError(err => {
        console.log(err);
        this.alerts.setAlert(`Failed to download v${pkg.version}`);
        this.log.error(`failed to download v${pkg.version}`, err);
        return throwError(err);
      }),
    );
  }

}
