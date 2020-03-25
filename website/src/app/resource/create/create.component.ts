import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup} from '@angular/forms';
import {ResourceService} from '../resource.service';
import {catchError, tap} from 'rxjs/operators';
import {throwError} from 'rxjs';
import {AlertService} from '../../alert.service';
import {HttpErrorResponse} from '@angular/common/http';
import {Router} from '@angular/router';

interface ResourceInputControls {
  name: string;
}

@Component({
  selector: 'app-resource-create',
  templateUrl: './create.component.html',
  styleUrls: ['./create.component.scss']
})
export class ResourceCreateComponent implements OnInit {

  form: FormGroup;

  constructor(
    private resources: ResourceService,
    private formBuilder: FormBuilder,
    private alerts: AlertService,
    private router: Router,
  ) {
    this.form = this.formBuilder.group({
      name: '',
    } as ResourceInputControls);
  }

  ngOnInit(): void {
  }

  onSubmit(input: ResourceInputControls) {
    this.form.disable();

    this.resources.create(input.name).pipe(
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
      })
    ).subscribe(data => {
      this.form.enable();
      this.alerts.setAlert(JSON.stringify(data));
      return this.router.navigate(['r', input.name]);
    });
  }

}
