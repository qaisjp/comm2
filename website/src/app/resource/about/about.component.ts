import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';

@Component({
  selector: 'app-about',
  templateUrl: './about.component.html',
  styleUrls: ['./about.component.scss']
})
export class ResourceAboutComponent implements OnInit {

  constructor(
    public view: ResourceViewService,
  ) { }

  ngOnInit(): void {
  }

}
