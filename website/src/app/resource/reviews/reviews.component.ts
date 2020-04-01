import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';

@Component({
  selector: 'app-reviews',
  templateUrl: './reviews.component.html',
  styleUrls: ['./reviews.component.scss']
})
export class ResourceReviewsComponent implements OnInit {

  constructor(
    public view: ResourceViewService,
  ) { }

  ngOnInit(): void {
  }

}
