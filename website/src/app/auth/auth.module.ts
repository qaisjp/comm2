import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login.component';
import { ReactiveFormsModule } from '@angular/forms';
import { LoginPageComponent } from './login-page/login-page.component';



@NgModule({
  declarations: [LoginComponent, LoginPageComponent],
  imports: [
    CommonModule,
    ReactiveFormsModule,
  ],
  exports: [
    LoginComponent
  ]
})
export class AuthModule { }
