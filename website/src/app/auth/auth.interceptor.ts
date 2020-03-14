import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler, HttpEvent } from '@angular/common/http';
import { AuthService } from './auth.service';
import { Observable } from 'rxjs';
import {environment} from '../../environments/environment';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private auth: AuthService) { }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {

    // Ensure the url starts with our API endpoint
    if (!req.url.startsWith(environment.api.baseurl)) {
      return;
    }

    // HACK
    if (req.headers.has('X-Authorization-None')) {
      const headers = req.headers.delete('X-Authorization-None');
      const directRequest = req.clone({ headers });
      return next.handle(directRequest);
    }

    if (this.auth.accessToken !== null) {
      req = req.clone({
        setHeaders: {
          Authorization: `Bearer ${this.auth.accessToken}`
        }
      });
    }

    return next.handle(req);
  }
}
