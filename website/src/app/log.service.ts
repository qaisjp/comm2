import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class LogService {
  log(...msg: any) {
    console.log(...msg);
  }

  error(...msg: any) {
    console.error(...msg);
  }

  warn(...msg: any) {
    console.warn(...msg);
  }

  debug(...msg: any) {
    console.debug(...msg);
  }
}
