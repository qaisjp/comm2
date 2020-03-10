import {Observable, of} from 'rxjs';

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
