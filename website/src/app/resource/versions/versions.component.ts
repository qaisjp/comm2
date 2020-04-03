import { Component, OnInit } from '@angular/core';
import {ResourceViewService} from '../resource-view.service';
import {Resource, ResourcePackage} from '../resource.service';

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

  download(res: Resource, pkg: ResourcePackage, anchor: HTMLAnchorElement) {
    this.view.download(pkg).subscribe((blob: Blob) => {
      const url = URL.createObjectURL(blob);
      console.log('we have', pkg.version, blob);

      anchor.href = url;
      anchor.download = res.name + '.zip';
      anchor.click();
      anchor.href = '';
    });
  }

}
