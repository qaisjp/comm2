import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { LogService } from '../log.service';
import { tap, catchError } from 'rxjs/operators';
import {Observable, of, ReplaySubject, Subject, throwError} from 'rxjs';
import { AuthenticatedUser } from '../user/user.service';

interface LoginResponse {
  token: string;
  expire: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(
    private http: HttpClient,
    private log: LogService
  ) { }

  accessToken: string | null = null;

  private userSource = new ReplaySubject<AuthenticatedUser>(1);
  user$ = this.userSource.asObservable();

  private sessionRestored = false;

  private static setAccessToken(token: string) {
    localStorage.setItem('accessToken', token);
  }

  public static canRestoreSession(): boolean {
    return localStorage.getItem('accessToken') !== null;
  }

  restoreSession(): Promise<AuthenticatedUser> {
    // Prevent multiple calls of restoreSession (without a manual clear)
    if (this.sessionRestored) {
      throw new Error('restoreSession should not be called multiple times');
    }

    const accessToken = localStorage.getItem('accessToken');
    if (accessToken === null) {
      this.log.log('auth.service/restoreSession: you are not logged in.');
      return throwError('you are not logged in').toPromise();
    }

    // Mark the session as restored before any request, to prevent concurrent restores
    this.sessionRestored = true;

    // Load access token from localStorage
    this.accessToken = localStorage.getItem('accessToken');

    // Prepare observable
    const success = new Subject<AuthenticatedUser>()

    // Get the local user
    this.http.get(`${environment.api.baseurl}/v1/user`).pipe(
      tap(userData => this.log.debug(`login get user response`, userData)),
      catchError(err => {
        this.logout({silent: true});

        success.error(err);
        // this.handleError<string>('login.get-user');
        return throwError(err);
      }),
    ).subscribe((userData: AuthenticatedUser) => {
      this.userSource.next(userData);
      success.next(userData);
    });

    return success.toPromise();
  }


  public login(username: string, password: string) {

    this.http.post(`${environment.api.baseurl}/v1/auth/login`, { username, password }).pipe(
      tap(data => this.log.debug(`login response`, data)),
        catchError(this.handleError<string>('login'))
    ).subscribe((data: LoginResponse) => {
      this.log.log('login response: ', data);

      AuthService.setAccessToken(data.token);
      this.sessionRestored = false;
      return this.restoreSession();
    });
  }

  public logout(config: { silent: boolean } = { silent: false }) {
    this.accessToken = null;
    localStorage.removeItem('accessToken');
    this.userSource.next(null);
  }

  /**
   *
   * Handle Http operation that failed.
   * Let the app continue.
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T>(operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
      // TODO: better job of transforming error for user consumption
      // TODO: send the error to remote logging infrastructure
      alert(`auth.service | ${operation} failed: ${error.message}`);
      this.log.error(error);

      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }
}
