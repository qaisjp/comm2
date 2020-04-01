import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';

@Component({
  selector: 'app-versions',
  templateUrl: './versions.component.html',
  styleUrls: ['./versions.component.scss']
})
export class ResourceVersionsComponent implements OnInit {

  constructor(
    public view: ResourceViewService,
  ) { }

  ngOnInit(): void {
  }

}
