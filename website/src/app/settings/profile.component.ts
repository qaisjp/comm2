import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup} from '@angular/forms';
import {UserProfileData, UserService} from '../user/user.service';
import {catchError, single, switchMap} from 'rxjs/operators';
import {Observable, throwError} from 'rxjs';
import {AlertService} from '../alert.service';
import {AuthService} from '../auth/auth.service';
import {Router} from '@angular/router';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class SettingsProfileComponent implements OnInit {
  form: FormGroup = this.formBuilder.group({
      bio: '',
      location: '',
      organisation: '',
      website: '',
    } as UserProfileData);

  constructor(
    private formBuilder: FormBuilder,
    private users: UserService,
    private alerts: AlertService,
    private auth: AuthService,
    private router: Router,
  ) { }

  ngOnInit(): void {
    this.form.disable();
    this.users.getProfile().pipe(
      catchError(e => {
        this.alerts.setAlert('Failed to load profile data');
        return throwError(e);
      }),
      single(),
    ).subscribe(data => {
      this.form.setValue(data);
      this.form.enable();
    });
  }

  onSubmit(data: UserProfileData) {
    this.users.patchProfile(data).pipe(
      catchError((reason: string) => {
        this.alerts.setAlert(reason);
        return throwError(reason);
      }),
      switchMap(() => this.auth.user$),
    ).subscribe(user => {
      this.router.navigate(['/u', user.username]);
    });
  }

}
