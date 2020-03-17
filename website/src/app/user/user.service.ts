import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {LogService} from '../log.service';
import {Observable} from 'rxjs';
import {environment} from '../../environments/environment';
import {catchError, map, tap} from 'rxjs/operators';
import {alertErrorReturnZero} from '../util';
import {Resource} from '../resource/resource.service';

export interface User {
  readonly id: number;
  readonly created_at: Date;
  username: string;
}

export interface AuthenticatedUser extends User {
  readonly updated_at: Date;
}

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor(
    private http: HttpClient,
    private log: LogService,
  ) {
  }

  public getUser(usernameOrID: string | number): Observable<User> {
    const url = `${environment.api.baseurl}/v1/users/${encodeURIComponent(usernameOrID)}`;
    return this.http.get(url, {headers: {'X-Authorization-None': ''}}).pipe(
      tap(data => this.log.debug(`getUser response`, data)),
      map(data => data as User)
    );
  }
}
