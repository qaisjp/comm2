import { Component, OnInit } from '@angular/core';
import {Resource, ResourceService} from '../resource/resource.service';
import {Observable} from 'rxjs';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit {
  latestResources$: Observable<Resource[]>;

  constructor(
    private resources: ResourceService
  ) {

  }

  ngOnInit() {
    this.latestResources$ = this.resources.getLatestResources();
  }

}
