import {Component, OnInit} from '@angular/core';
import {AuthService} from './auth/auth.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = 'mtahub';

  constructor(
    private auth: AuthService,
  ) { }

  ngOnInit() {
    if (AuthService.canRestoreSession()) {
      this.auth.restoreSession().catch(reason => {
        console.error('restoreSession on initial start failed because', reason);
      });
    }
  }
}
