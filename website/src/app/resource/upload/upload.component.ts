import { Component, OnInit } from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {ResourcePackage, ResourceService} from '../resource.service';
import {ResourceViewService} from '../resource-view.service';
import {AlertService} from '../../alert.service';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.scss']
})
export class ResourceUploadComponent implements OnInit {
  editMode: boolean;
  pkg: ResourcePackage;

  constructor(
    private route: ActivatedRoute,
    private resources: ResourceService,
    private view: ResourceViewService,
    private alerts: AlertService,
  ) {

  }

  ngOnInit(): void {
    this.route.data.subscribe(({ pkg }: { pkg: ResourcePackage }) => {
      this.editMode = pkg !== undefined;
      if (!this.editMode) {
        return;
      }
      this.pkg = pkg;
      console.log(pkg);
    });
  }

}
