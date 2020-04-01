import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';

@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ResourceManageComponent implements OnInit {

  constructor(
    public view: ResourceViewService,
  ) { }

  ngOnInit(): void {
  }

}
