import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {ResourceStatus} from '../resource/resource.service';
import {User, UserProfile, UserService} from '../user/user.service';
import {AlertService} from '../alert.service';
import {Location} from '@angular/common';
import {Subject} from 'rxjs';

interface UserProfileExtended extends UserProfile {
  hasPrivate: boolean;
}

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss']
})
export class ProfileComponent implements OnInit {
  public user$ = new Subject<UserProfileExtended>();

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
      });
    });
  }

}
