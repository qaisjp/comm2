import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {LogService} from '../log.service';
import {Observable} from 'rxjs';
import {environment} from '../../environments/environment';
import {catchError, map, tap} from 'rxjs/operators';
import {alertErrorReturnZero} from '../util';

export enum ResourceStatus {
  PUBLIC = 'public',
  PRIVATE = 'private',
}

export interface Resource {
  readonly id: number;
  readonly created_at: Date;
  readonly updated_at: Date;
  author_id: number;
  name: string;
  title: string;
  description: string;
  status: ResourceStatus;
}

interface ResourceCreateResponse {
  readonly id: number;
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

  public get(id: number | string): Observable<Resource> {
    return this.http.get(`${environment.api.baseurl}/v1/resources/${encodeURIComponent(id)}`).pipe(
      tap(data => this.log.debug(`getResource(${id})`)),
      map(data => data as Resource),
    );
  }

  public getLatestResources(): Observable<Resource> {
    return this.http.get(`${environment.api.baseurl}/v1/resources`, {headers: {'X-Authorization-None': ''}}).pipe(
      tap(data => this.log.debug(`getLatestResources response`, data)),
      catchError(alertErrorReturnZero<string>('ResourceService.getLatestResources')),
      map(data => data as Resource),
    );
  }

  public create(name: string): Observable<ResourceCreateResponse> {
    return this.http.post(`${environment.api.baseurl}/v1/resources`, {name}).pipe(
      tap(data => this.log.debug(`sending createResource with name ${name}`)),
      map(data => data as ResourceCreateResponse),
    );
  }
}
