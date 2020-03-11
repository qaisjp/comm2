import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import { LoginPageComponent } from './login-page/login-page.component';
import {RouterModule} from '@angular/router';
import { RegisterPageComponent } from './register-page/register-page.component';



@NgModule({
  declarations: [LoginComponent, LoginPageComponent, RegisterPageComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    RouterModule,
  ],
  exports: [
    LoginComponent
  ]
})
export class AuthModule { }
