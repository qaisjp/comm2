import {Observable, of} from 'rxjs';
import {HttpErrorResponse} from '@angular/common/http';
import {CONFLICT as HTTP_STATUS_CONFLICT} from 'http-status-codes';

export function alertErrorReturnZero<T>(operation = 'operation', result?: T) {
  return (error: any): Observable<T> => {
    // TODO: better job of transforming error for user consumption
    // TODO: send the error to remote logging infrastructure
    alert(`${operation} failed: ${error.message}`);
    this.log.error(error);

    // Let the app keep running by returning an empty result.
    return of(result as T);
  };
}

export function httpErrToMessage(err: HttpErrorResponse) {
  let reason = 'Something went wrong.';
  if (err.ok) {
    reason = 'Account already exists with that username or email address.';
  }
}
