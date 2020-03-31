import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup} from '@angular/forms';
import {AuthService} from '../auth.service';
import {Router} from '@angular/router';
import {AlertService} from '../../alert.service';
import {LogService} from '../../log.service';
import {throwError} from 'rxjs';
import {catchError} from 'rxjs/operators';
import {HttpErrorResponse} from '@angular/common/http';
import {CONFLICT as HTTP_STATUS_CONFLICT} from 'http-status-codes';

interface RegisterInputControls {
  username: string;
  email: string;
  password: string;
}

@Component({
  selector: 'app-register-page',
  templateUrl: './register-page.component.html',
  styleUrls: ['./register-page.component.scss']
})
export class RegisterPageComponent implements OnInit {
  form: FormGroup;

  constructor(
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private router: Router,
    private alerts: AlertService,
    private log: LogService,
  ) {
    this.form = this.formBuilder.group({
      username: '',
      email: '',
      password: '',
    } as RegisterInputControls);
  }

  ngOnInit() {
  }

  onSubmit(data: RegisterInputControls) {
    this.form.disable();
    this.log.debug('register form data', data);
    this.authService.register(data.username, data.email, data.password).pipe(
      catchError((httpErrorResponse: HttpErrorResponse) => {
        let reason = 'Something went wrong.';
        if (httpErrorResponse.status === HTTP_STATUS_CONFLICT) {
          reason = 'Account already exists with that username or email address.';
        }
        this.alerts.setAlert(reason);
        this.form.enable();
        return throwError(reason);
      })
    ).subscribe(() => {
      this.alerts.clearAlert();
      return this.router.navigate(['/login']);
    });
  }

}
