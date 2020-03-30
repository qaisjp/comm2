import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {ResourceStatus} from '../resource/resource.service';
import {User, UserProfile, UserService} from '../user/user.service';
import {AlertService} from '../alert.service';
import {Location} from '@angular/common';
import {Observable, of, ReplaySubject, Subject} from 'rxjs';
import {AuthService} from '../auth/auth.service';
import {delay, single, switchMap} from 'rxjs/operators';

interface UserProfileExtended extends UserProfile {
  hasPrivate: boolean;
}

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  public user$ = new ReplaySubject<UserProfileExtended>(1);
  public followed = false;
  public loading = false; // HACK

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private location: Location,
    private users: UserService,
    private alerts: AlertService,
    public auth: AuthService,
  ) {
  }

  ngOnInit() {
    this.route.params.subscribe(params => {
      this.users.getUserProfile(params.username).subscribe((data: UserProfile) => {
        // Update url from ID to username if necessary without causing a page reload
        if (data.username !== params.username) {
          this.router.navigate(['u', data.username], {
            preserveFragment: true,
            queryParamsHandling: 'preserve',
            replaceUrl: true,
          });
        }

        const hasPrivate = undefined !==
          data.resources.find(r => r.status === ResourceStatus.PRIVATE);

        this.user$.next({
          ...data,
          hasPrivate,
        });

        this.auth.user$.subscribe(
          user => this.followed = data.followers.some(u => u.id === user.id)
        );
      });
    });
  }

  toggleFollowState() {
    // are we expected to subscribe in every single function? observables are annoying
    this.user$.subscribe(user => {
      this.loading = true;
      let obs;
      if (this.followed) {
        obs = this.users.unfollowUser(user.id);
      } else {
        obs = this.users.followUser(user.id);
      }
      obs.subscribe(() => {
        this.followed = !this.followed;
        this.loading = false;
      });
    });
  }

}
