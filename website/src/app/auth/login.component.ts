import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { AuthService } from './auth.service';

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
  ) {
    this.loginForm = this.formBuilder.group({
      username: '',
      password: '',
    } as LoginInputControls)
  }

  ngOnInit() {
  }

  onSubmit(data: LoginInputControls) {
    this.loginForm.disable();
    console.log("login form data", data);
    this.authService.login(data.username, data.password);
  }

}
