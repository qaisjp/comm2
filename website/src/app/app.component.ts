import {Component, OnInit} from '@angular/core';
import {AuthService} from './auth/auth.service';
import {AlertService} from './alert.service';
import {Router} from '@angular/router';
import {catchError} from 'rxjs/operators';
import {throwError} from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'mtahub';

  constructor(
    private auth: AuthService,
    public alerts: AlertService,
  ) {
  }

  ngOnInit() {
    if (!AuthService.canRestoreSession()) {
      return;
    }

    this.auth.restoreSession().pipe(
      catchError(reason => {
        console.error('restoreSession on initial start failed because', reason);
        this.alerts.setAlert(reason.message); // todo actually fix reason
        return throwError(reason);
      })
    ).subscribe(u => {
      console.log('Logged in', u.username);
    });

  }
}
