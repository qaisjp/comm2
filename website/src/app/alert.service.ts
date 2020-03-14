import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AlertService {

  private alertSource = new Subject<string>();
  alert$ = this.alertSource.asObservable();

  constructor() { }

  public setAlert(str: string) {
    this.alertSource.next(str);
  }

  public clearAlert() {
    this.alertSource.next(null);
  }
}
