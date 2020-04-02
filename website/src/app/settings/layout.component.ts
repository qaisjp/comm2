import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from '@angular/router';

@Component({
  selector: 'app-account',
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.scss']
})
export class SettingsLayoutComponent implements OnInit {

  constructor(
    public route: ActivatedRoute,
  ) { }

  ngOnInit(): void {
  }

}
