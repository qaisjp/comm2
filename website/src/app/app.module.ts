import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { ProfileModule } from './profile/profile.module';
import { AuthModule } from './auth/auth.module';
import { AuthInterceptor } from './auth/auth.interceptor';
import { LogService } from './log.service';
import { NavbarComponent } from './navbar/navbar.component';
import {ResourceModule} from './resource/resource.module';
import { AdminModule } from './admin/admin.module';
import {OcticonModule} from './octicon/octicon.module';
import { SettingsProfileComponent } from './settings/profile.component';
import { SettingsAccountComponent } from './settings/account.component';
import {ReactiveFormsModule} from '@angular/forms';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    NavbarComponent,
    SettingsProfileComponent,
    SettingsAccountComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule, // import HttpClientModule after BrowserModule
    AppRoutingModule,
    OcticonModule,
    ResourceModule,
    AuthModule,
    ReactiveFormsModule,
  ],
  providers: [
    { provide: HTTP_INTERCEPTORS, useClass: AuthInterceptor, multi: true },
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
