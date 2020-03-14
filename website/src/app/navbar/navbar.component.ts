import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth/auth.service';
import {AuthenticatedUser} from '../user/user.service';
import {Observable} from 'rxjs';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
  private user$: Observable<AuthenticatedUser>;

  constructor(
    private auth: AuthService) {
    this.user$ = auth.user$;
  }

  ngOnInit() {
  }

}
