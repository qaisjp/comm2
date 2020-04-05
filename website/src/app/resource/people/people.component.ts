import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';
import {FormBuilder, FormGroup} from '@angular/forms';

@Component({
  selector: 'app-people',
  templateUrl: './people.component.html',
  styleUrls: ['./people.component.scss']
})
export class ResourcePeopleComponent implements OnInit {

  form: FormGroup = this.fb.group({username:''});

  constructor(
    public view: ResourceViewService,
    private fb: FormBuilder,
  ) { }

  ngOnInit(): void {
  }

  remove(uid: number) {

  }

  add(username) {

  }
}
