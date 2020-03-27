import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup} from '@angular/forms';
import {ResourceService} from '../resource.service';
import {catchError, switchMap, tap} from 'rxjs/operators';
import {throwError} from 'rxjs';
import {AlertService} from '../../alert.service';
import {HttpErrorResponse} from '@angular/common/http';
import {Router} from '@angular/router';
import {AuthService} from '../../auth/auth.service';
import {AuthenticatedUser} from '../../user/user.service';

interface ResourceInputControls {
  name: string;
  title: string;
  description: string;
}

@Component({
  selector: 'app-resource-create',
  templateUrl: './create.component.html',
  styleUrls: ['./create.component.scss']
})
export class ResourceCreateComponent implements OnInit {

  form: FormGroup;
  titleInput = false;

  constructor(
    private resources: ResourceService,
    private formBuilder: FormBuilder,
    private alerts: AlertService,
    private auth: AuthService,
    private router: Router,
  ) {
    this.form = this.formBuilder.group({
      name: '',
      title: '',
      description: '',
    } as ResourceInputControls);
  }

  ngOnInit(): void {
  }

  onSubmit(input: ResourceInputControls) {
    this.form.disable();

    this.resources.create(input.name, input.title, input.description).pipe(
      catchError((err: HttpErrorResponse) => {
        this.form.enable();
        let msg: string;
        if (err.error instanceof ErrorEvent) {
          msg = `An error occurred: ${err.error.message}`;
        } else {
          msg = err.error.message;
        }
        this.alerts.setAlert(msg);
        return throwError(msg);
      }),
      switchMap(data => {
        this.form.enable();
        this.alerts.setAlert(JSON.stringify(data));
        return this.auth.user$;
      }),
      switchMap((user: AuthenticatedUser) => this.router.navigate(['/u', user.username, input.name])),
    ).subscribe();
  }

  updateTitle(title: string) {
    // Don't update the title manually if the user has written stuff there
    if (this.titleInput) {
      return;
    }

    this.form.patchValue({title});
  }

}
