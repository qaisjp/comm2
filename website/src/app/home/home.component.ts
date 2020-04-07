import { Component, OnInit } from '@angular/core';
import {Resource, ResourceService} from '../resource/resource.service';
import {NEVER, Observable} from 'rxjs';
import {environment} from '../../environments/environment';
import {catchError, map, take, tap} from 'rxjs/operators';
import {HttpClient, HttpErrorResponse} from '@angular/common/http';
import {LogService} from '../log.service';
import {AlertService} from '../alert.service';
import {User} from '../user/user.service';

interface LatestResource extends Resource {
  author_username: User['username'];
}

interface Homepage {
  latest: LatestResource[];
}

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit {
  data?: Homepage = null;

  constructor(
    private resources: ResourceService,
    private http: HttpClient,
    private log: LogService,
    private alerts: AlertService,
  ) {

  }

  ngOnInit() {
    const url = `${environment.api.baseurl}/private/homepage`;
    this.http.get(url, {headers: {'X-Authorization-None': ''}}).pipe(
      take(1),
      tap(data => this.log.debug(`getHomepage response`, data)),
      catchError((err: HttpErrorResponse) => {
        this.alerts.setAlert(err.message);
        return NEVER;
      }),
      map(data => data as Homepage),
    ).subscribe(data => {
      this.data = data;
    });
  }

}
