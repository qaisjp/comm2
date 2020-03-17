import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { AuthService } from './auth.service';
import {AlertService} from '../alert.service';
import {Router} from '@angular/router';
import {catchError} from 'rxjs/operators';
import {throwError} from 'rxjs';

interface LoginInputControls {
  username: string,
  password: string,
}

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  loginForm: FormGroup;

  constructor(
    private authService: AuthService,
    private formBuilder: FormBuilder,
    private router: Router,
    private alerts: AlertService,
  ) {
    this.loginForm = this.formBuilder.group({
      username: '',
      password: '',
    } as LoginInputControls);
  }

  ngOnInit() {
  }

  onSubmit(data: LoginInputControls) {
    this.loginForm.disable();
    console.log('login form data', data);


    this.authService.login(data.username, data.password).pipe(
      catchError(reason => {
        this.alerts.setAlert(reason);
        this.loginForm.enable();
        return throwError(reason);
      })
    ).subscribe(() => {
      return this.router.navigate(['/']);
    });

  }

}
