import { Injectable } from '@angular/core';
import {HttpClient, HttpErrorResponse, HttpEvent, HttpRequest} from '@angular/common/http';
import {LogService} from '../log.service';
import {Observable, throwError} from 'rxjs';
import {environment} from '../../environments/environment';
import {catchError, map, tap} from 'rxjs/operators';
import {alertErrorReturnZero} from '../util';
import {User, UserID} from '../user/user.service';
import {BAD_REQUEST} from 'http-status-codes';

export enum ResourceVisibility {
  PUBLIC = 'public',
  PRIVATE = 'private',
}

export interface Resource {
  readonly id: number;
  readonly created_at: string;
  readonly updated_at: string;
  author_id: number;
  name: string;
  title: string;
  description: string;
  visibility: ResourceVisibility;
  archived: boolean;
  authors: User[];
  readonly can_manage: boolean;
}

export type ResourceCreateResponse = Readonly<Pick<Resource, 'id'>>;
export type ResourcePatchRequest = Partial<Pick<Resource, 'name' | 'title' | 'description' | 'visibility' | 'archived'>>;

// ResourceID can either be the name of the resource, or its ID
export type ResourceID = Resource['id'] | Resource['name'];

export type PackageID = number;

export interface ResourcePackage {
  readonly id: number;
  readonly created_at: string;
  readonly updated_at: string;
  published_at?: string;

  readonly resource_id: number;
  readonly author_id: number;
  version: string;
  description: string;
  file_uploaded: boolean;
}

@Injectable({
  providedIn: 'root'
})
export class ResourceService {
  constructor(
    private http: HttpClient,
    private log: LogService,
  ) {
  }

  private getResourceURL(userID: UserID, resourceID: ResourceID): string {
    return `${environment.api.baseurl}/v1/resources/${encodeURIComponent(userID)}/${encodeURIComponent(resourceID)}`;
  }

  public get(userID: UserID, resourceID: ResourceID): Observable<Resource> {
    return this.http.get(this.getResourceURL(userID, resourceID)).pipe(
      tap(data => this.log.debug(`getResource(${userID}/${resourceID})`)),
      map(data => data as Resource),
    );
  }

  public getLatest(): Observable<Resource[]> {
    return this.http.get(`${environment.api.baseurl}/v1/resources`, {headers: {'X-Authorization-None': ''}}).pipe(
      tap(data => this.log.debug(`getLatestResources response`, data)),
      catchError(alertErrorReturnZero<string>('ResourceService.getLatestResources')),
      map(data => data as Resource[]),
    );
  }

  public create(name: string, title: string, description: string): Observable<ResourceCreateResponse> {
    return this.http.post(`${environment.api.baseurl}/v1/resources`, {name, title, description}).pipe(
      tap(data => this.log.debug(`sending createResource with ${JSON.stringify({name, title, description})}`)),
      map(data => data as ResourceCreateResponse),
    );
  }

  public patch(userID: UserID, resourceID: ResourceID, reqData: ResourcePatchRequest): Observable<void> {
    return this.http.patch(this.getResourceURL(userID, resourceID), reqData).pipe(
        tap(data => this.log.debug(`patch(${userID}, ${resourceID}) with req ${JSON.stringify(reqData)}`)),
        catchError((err: HttpErrorResponse) => {
          let reason = 'Something went wrong';
          if (err.status === BAD_REQUEST) {
            reason = err.error.message;
          }
          return throwError(reason);
        }),
        map(() => void 0),
    );
  }

  public getPackages(userID: UserID, resourceID: ResourceID): Observable<ResourcePackage[]> {
    return this.http.get(`${this.getResourceURL(userID, resourceID)}/pkg`).pipe(
      tap(data => this.log.debug(`getResourcePackages(${userID}, ${resourceID}`)),
      map(data => data as ResourcePackage[])
    );
  }

  public getPackage(userID: UserID, resourceID: ResourceID, packageID: PackageID): Observable<ResourcePackage> {
    return this.http.get(`${this.getResourceURL(userID, resourceID)}/pkg/${encodeURIComponent(packageID)}`).pipe(
      tap(data => this.log.debug(`getResourcePackage(${userID}, ${resourceID}, ${packageID}`)),
      map(data => data as ResourcePackage)
    );
  }

  transfer(userID: UserID, resourceID: ResourceID, username: string) {
    return this.http.post(this.getResourceURL(userID, resourceID) + '/transfer', {new_owner: username}).pipe(
        tap(data => this.log.debug(`transferResource(${userID}, ${resourceID}) to username ${username}`)),
        catchError((err: HttpErrorResponse) => {
          let reason = 'Something went wrong';
          if (err.status === BAD_REQUEST) {
            reason = err.error.message;
          }
          return throwError(reason);
        }),
        map(data => data as {new_username: string}),
    );
  }

  delete(userID: UserID, resourceID: ResourceID) {
    return this.http.delete(this.getResourceURL(userID, resourceID)).pipe(
        tap(data => this.log.debug(`deleteResource(${userID}, ${resourceID})`)),
        catchError((err: HttpErrorResponse) => {
          let reason = 'Something went wrong';
          if (err.status === BAD_REQUEST) {
            reason = err.error.message;
          }
          return throwError(reason);
        }),
        map(data => void 0),
    );
  }

  download(userID: UserID, resourceID: ResourceID, packageID: PackageID): Observable<HttpEvent<any>> {
    const url = this.getResourceURL(userID, resourceID) + `/pkg/${encodeURIComponent(packageID)}/download`;
    const req = new HttpRequest('GET', url, {reportProgress: true, responseType: 'blob'});
    return this.http.request(req).pipe(
      tap(data => this.log.debug(`downloadResource(${userID}, ${resourceID}, ${packageID})`))
    );
  }
}
