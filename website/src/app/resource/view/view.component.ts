import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Resource, ResourcePackage, ResourceService} from '../resource.service';
import {AlertService} from '../../alert.service';
import {pluck, switchMap, tap} from 'rxjs/operators';
import {Observable, Subject} from 'rxjs';

@Component({
  selector: 'app-resource-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss']
})
export class ResourceViewComponent implements OnInit {
  public resource$ = new Subject<Resource>();
  public packages$: Observable<ResourcePackage[]>;

  constructor(
    private route: ActivatedRoute,
    private resources: ResourceService,
    private alerts: AlertService,
  ) { }

  ngOnInit() {
    this.route.params.pipe(
      switchMap(params => this.resources.get(params.username, params.resource))
    ).subscribe((data: Resource) => {
      this.resource$.next(data);
      this.packages$ = this.resources.getPackages(data.author_id, data.id);
    });
  }

}
