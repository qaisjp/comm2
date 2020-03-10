import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment';
import { LogService } from '../log.service';
import { tap, catchError } from 'rxjs/operators';
import { Observable, of, ReplaySubject } from 'rxjs';
import { AuthenticatedUser } from '../user/user.service';

interface LoginResponse {
  token: string;
  expire: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  accessToken: string | null = null;

  private userSource = new ReplaySubject<AuthenticatedUser>(1);
  user$ = this.userSource.asObservable();

  constructor(
    private http: HttpClient,
    private log: LogService
  ) { }

  public login(username: string, password: string) {
    this.http.post(`${environment.api.baseurl}/v1/auth/login`, { username, password }).pipe(
      tap(data => this.log.debug(`login response`, data)),
        catchError(this.handleError<string>('login'))
    ).subscribe((data: LoginResponse) => {
      this.log.log('login response: ', data);
      this.accessToken = data.token;

      this.http.get(`${environment.api.baseurl}/v1/user`).pipe(
        tap(userData => this.log.debug(`login get user response`, userData)),
        catchError(this.handleError<string>('login.get-user'))
      ).subscribe((userData: AuthenticatedUser) => {
        this.userSource.next(userData);
      });
    });
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
