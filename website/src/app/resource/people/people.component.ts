import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';

@Component({
  selector: 'app-people',
  templateUrl: './people.component.html',
  styleUrls: ['./people.component.scss']
})
export class ResourcePeopleComponent implements OnInit {

  constructor(
    public view: ResourceViewService,
  ) { }

  ngOnInit(): void {
  }

}
