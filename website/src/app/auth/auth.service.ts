import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private http: HttpClient) { }

  public login(username: string, password: string) {
    this.http.post("http://localhost:8080/v1/auth/login", { username, password })
      .subscribe((data) => {
        console.log(data);
      })
  }
}
