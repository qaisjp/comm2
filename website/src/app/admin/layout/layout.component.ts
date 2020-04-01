import { Component, OnInit } from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
  selector: 'app-layout',
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.scss'],
  preserveWhitespaces: true,
})
export class AdminLayoutComponent implements OnInit {

  constructor(
    public route: ActivatedRoute,
    public router: Router,
  ) { }

  ngOnInit(): void {

  }

}
