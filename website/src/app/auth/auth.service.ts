import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/environments/environment'

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  accessToken: string | null = null;

  constructor(private http: HttpClient) { }

  public login(username: string, password: string) {
    this.http.post(`${environment.api.baseurl}/v1/auth/login`, { username, password })
      .subscribe((data) => {
        console.log(data);
      })
  }
}
