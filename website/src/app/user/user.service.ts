import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse} from '@angular/common/http';
import {LogService} from '../log.service';
import {Observable, throwError} from 'rxjs';
import {environment} from '../../environments/environment';
import {catchError, map, tap} from 'rxjs/operators';
import {alertErrorReturnZero} from '../util';
import {Resource, ResourceID} from '../resource/resource.service';
import {UNAUTHORIZED as HTTP_STATUS_UNAUTHORIZED, CONFLICT as HTTP_STATUS_CONFLICT} from 'http-status-codes';

export interface User {
  readonly id: number;
  readonly created_at: string;
  readonly username: string;
  readonly gravatar: string;
  readonly level: number;

  readonly follows_you?: boolean;
}

export interface AuthenticatedUser extends User {
  readonly updated_at: string;
}

export interface UserProfileData {
  location: string;
  organisation: string;
  website: string;
  bio: string;
}

export interface UserProfile extends User, UserProfileData {
  readonly resources: Resource[];
  readonly following: User[];
  readonly followers: User[];
}

// UserID can either be the name of the user, or its ID
export type UserID = string | number;

@Injectable({
  providedIn: 'root'
})
export class UserService {

  constructor(
    private http: HttpClient,
    private log: LogService,
  ) {
  }

  // CURRENT USER ONLY
  // todo: probably move to AuthenticatedUserService (?)
  delete() {
    return this.http.delete(`${environment.api.baseurl}/private/account`).pipe(
      tap(data => this.log.debug(`deleteMyAccount()`)),
      map(data => void 0),
    );
  }

  chgPass(password: string, newPassword: string) {
    const body = {password, new_password: newPassword};
    return this.http.post( `${environment.api.baseurl}/private/account/password`, body).pipe(
      tap(data => this.log.debug(`changePassword()`)),
      catchError((err: HttpErrorResponse) => {
        let reason = 'Something went wrong';
        if (err.status === HTTP_STATUS_UNAUTHORIZED) {
          reason = err.error.message;
        }
        return throwError(reason);
      }),
      map(data => void 0),
    );
  }

  rename(username: string) {
    const body = {username};
    return this.http.post( `${environment.api.baseurl}/private/account/username`, body).pipe(
      tap(data => this.log.debug(`changeUsername(to=${username})`)),
      catchError((err: HttpErrorResponse) => {
        let reason = 'Something went wrong';
        if (err.status === HTTP_STATUS_CONFLICT) {
          reason = err.error.message;
        }
        return throwError(reason);
      }),
      map(data => void 0),
    );
  }

  // ALL USERS
  public getUser(id: UserID): Observable<User> {
    const url = `${environment.api.baseurl}/v1/users/${encodeURIComponent(id)}`;
    return this.http.get(url, {headers: {'X-Authorization-None': ''}}).pipe(
      tap(data => this.log.debug(`getUser response`, data)),
      map(data => data as User)
    );
  }

  public getUserProfile(id: UserID): Observable<UserProfile> {
    const url = `${environment.api.baseurl}/private/profile/${encodeURIComponent(id)}`;
    return this.http.get(url).pipe(
      tap(data => this.log.debug(`getUserProfile response`, data)),
      map(data => data as UserProfile)
    );
  }

  public followUser(id: UserID) {
    const url = `${environment.api.baseurl}/v1/user/follow/${id}`;
    return this.http.put(url, '').pipe(
      tap(data => this.log.debug(`delete user follow on`, id))
    );
  }

  public unfollowUser(id: UserID) {
    const url = `${environment.api.baseurl}/v1/user/follow/${id}`;
    return this.http.delete(url).pipe(
      tap(data => this.log.debug(`delete user follow on`, id))
    );
  }
}
