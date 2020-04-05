import {Component, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Resource, ResourcePackage, ResourceService} from '../resource.service';
import {AlertService} from '../../alert.service';
import {single, switchMap, tap} from 'rxjs/operators';
import {Observable, Subject} from 'rxjs';
import {ResourceViewService} from '../resource-view.service';

@Component({
  selector: 'app-resource-layout',
  templateUrl: './layout.component.html',
  styleUrls: ['./layout.component.scss']
})
export class ResourceLayoutComponent implements OnInit {
  constructor(
    public route: ActivatedRoute,
    private resources: ResourceService,
    private alerts: AlertService,
    public view: ResourceViewService,
  ) {
  }

  ngOnInit() {
    this.view.reinit();
    this.route.params.pipe(
      switchMap(params => this.resources.get(params.username, params.resource))
    ).subscribe(this.view.resource$);
  }

  downloadLatestPackage(res: Resource, anchor: HTMLAnchorElement) {
    this.view.packages$.pipe(
      single(),
    ).subscribe(packages => {
      const pkg = packages.filter(p => p.published_at)[0];
      this.view.download(pkg).subscribe((blob: Blob) => {
        // duplicated in versions.component.ts
        const url = URL.createObjectURL(blob);
        console.log('we have', pkg.version, blob);

        anchor.href = url;
        anchor.download = res.name + '.zip';
        anchor.click();
        anchor.href = '';
      });
  });
  }

}
