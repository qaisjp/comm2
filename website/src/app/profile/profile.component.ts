import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Params, Router} from '@angular/router';
import {mergeMap, pluck} from 'rxjs/operators';
import {ResourceService} from '../resource/resource.service';
import {User, UserService} from '../user/user.service';
import {AlertService} from '../alert.service';
import {Location} from '@angular/common';
import {Subject} from 'rxjs';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  public user$ = new Subject<User>();

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private location: Location,
    private users: UserService,
    private alerts: AlertService,
  ) {
  }

  ngOnInit() {
    this.route.params.subscribe(params => {
      this.users.getUser(params.username).subscribe((data: User) => {
        // Update url from ID to username if necessary without causing a page reload
        if (data.username !== params.username) {
          this.router.navigate(['u', data.username], {
            preserveFragment: true,
            queryParamsHandling: 'preserve',
            replaceUrl: true,
          });
        }

        this.user$.next(data);
        this.alerts.setAlert( JSON.stringify(data) );
      });
    });
  }

}
