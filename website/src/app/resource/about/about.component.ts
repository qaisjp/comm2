import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';
import {FormBuilder, FormGroup} from '@angular/forms';
import {catchError, map, single, switchMap, take, tap} from 'rxjs/operators';
import {ResourceService} from '../resource.service';
import {AlertService} from '../../alert.service';
import {NEVER, throwError} from 'rxjs';

@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class ResourceAboutComponent implements OnInit {
  // todo: needs to support cancel
  editing = false;
  form: FormGroup = this.fb.group({
    title: '',
    description: '',
  });


  constructor(
    public view: ResourceViewService,
    private resources: ResourceService,
    private fb: FormBuilder,
    private alerts: AlertService,
  ) { }

  ngOnInit(): void {
  }

  onLoad() {
    this.view.resource$.pipe(take(1)).subscribe(r => {
      this.editing = true;
      this.form.setValue({
        title: r.title,
        description: r.description,
      });
    });
  }

  onSave(data) {
    const reqData = {
      title: data.title,
      description: data.description,
    };
    this.form.disable();
    this.view.resource$.pipe(
      take(1),
      tap(_ => console.log("yep")),
      switchMap(r => this.resources.patch(r.author_id, r.id, reqData)),
      catchError(reason => {
        this.alerts.setAlert(reason);
        this.form.enable();
        return NEVER;
      }),
      switchMap(_ => this.view.resource$),
      take(1),
      map(r => ({
        ...r,
        title: reqData.title,
        description: reqData.description,
      })),
      tap(this.view.resource$),
    ).subscribe(r => {
      this.editing = false;
      this.form.enable();
    });
  }

}
