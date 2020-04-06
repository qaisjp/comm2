import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';
import {FormBuilder, FormGroup} from '@angular/forms';
import {ResourceService} from '../resource.service';
import {catchError, switchMap, take} from 'rxjs/operators';
import {AlertService} from '../../alert.service';
import {NEVER} from 'rxjs';

@Component({
  selector: 'app-people',
  templateUrl: './people.component.html',
  styleUrls: ['./people.component.scss']
})
export class ResourcePeopleComponent implements OnInit {

  form: FormGroup = this.fb.group({username:''});

  constructor(
    public view: ResourceViewService,
    private fb: FormBuilder,
    private resources: ResourceService,
    private alerts: AlertService,
  ) { }

  ngOnInit(): void {
  }

  remove(uid: number) {
    this.view.resource$.pipe(
      take(1),
      switchMap(r => this.resources.delCollab(r.author_id, r.id, uid)),
      catchError(reason => {
        this.alerts.setAlert(reason);
        return NEVER;
      })
    ).subscribe(_ => {
      this.view.refresh();
    });
  }

  add(username) {
    this.view.resource$.pipe(
      take(1),
      switchMap(r => this.resources.addCollab(r.author_id, r.id, username)),
      catchError(reason => {
        this.alerts.setAlert(reason);
        return NEVER;
      })
    ).subscribe(_ => {
      this.view.refresh();
    });
  }
}
