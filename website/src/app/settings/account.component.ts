import { Component, OnInit } from '@angular/core';
import {FormBuilder} from '@angular/forms';
import {Router} from '@angular/router';
import {UserService} from '../user/user.service';
import {catchError} from 'rxjs/operators';
import {throwError} from 'rxjs';
import {AlertService} from '../alert.service';
import {Location} from '@angular/common';

@Component({
  selector: 'app-account',
  templateUrl: './account.component.html',
  styleUrls: ['./account.component.scss']
})
export class SettingsAccountComponent implements OnInit {

  usernameForm = this.formBuilder.group({username: ''});
  deleteForm = this.formBuilder.group({});
  passwordForm = this.formBuilder.group({
    password: '',
    newPassword: '',
    newPasswordAgain: '',
  });

  constructor(
    private formBuilder: FormBuilder,
    private router: Router,
    private users: UserService,
    private alerts: AlertService,
  ) { }

  ngOnInit(): void {
  }

  submitChangeUsername({username}: {username: string}) {
    if (!confirm(`Are you sure you want to change username to ` + username)) {
      return;
    }

    this.users.rename(username).pipe(
      catchError((reason: string) => {
        this.alerts.setAlert(reason);
        return throwError(reason);
      }),
    ).subscribe(() => {
      this.router.navigate(['/u', username]);
    });
  }

  submitDelete(data) {
    if (!confirm(`Are you really sure you want to delete your account?`)) {
      return;
    }

    this.users.delete().pipe(
      catchError((reason: string) => {
        this.alerts.setAlert(reason);
        return throwError(reason);
      }),
    ).subscribe(() => {
      window.location.href = '/';
    });
  }

  submitChangePassword(data: {[key: string]: string}) {
    if (data.newPassword !== data.newPasswordAgain) {
      this.alerts.setAlert('Passwords don\'t match');
      return;
    }

    this.users.chgPass(data.password, data.newPassword).pipe(
      catchError((reason: string) => {
        this.alerts.setAlert(reason);
        return throwError(reason);
      }),
    ).subscribe(() => {
      this.alerts.setAlert('Success!');
    });
  }

}
