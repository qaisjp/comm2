import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Resource, ResourceService} from '../resource.service';
import {AlertService} from '../../alert.service';
import {pluck, switchMap} from 'rxjs/operators';
import {Subject} from 'rxjs';

@Component({
  selector: 'app-resource-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss']
})
export class ResourceViewComponent implements OnInit {
  public resource$ = new Subject<Resource>();

  constructor(
    private route: ActivatedRoute,
    private resources: ResourceService,
    private alerts: AlertService,
  ) { }

  ngOnInit() {
    this.route.params.pipe(
      pluck('resourceName'),
      switchMap(name => this.resources.get(name))
    ).subscribe((data: Resource) => {
      this.resource$.next(data);
    });
  }

}
