import {Component, OnInit} from '@angular/core';
import {ResourceViewService} from '../resource-view.service';
import {Resource, ResourcePackage, ResourceService, ResourceVisibility} from '../resource.service';
import {AlertService} from '../../alert.service';
import {catchError} from 'rxjs/operators';
import {throwError} from 'rxjs';
import {FormBuilder, FormGroup} from '@angular/forms';
import {Router} from '@angular/router';

@Component({
  selector: 'app-manage',
  templateUrl: './manage.component.html',
  styleUrls: ['./manage.component.scss']
})
export class ResourceManageComponent implements OnInit {
  resource: Resource = null;
  packages: ResourcePackage[] = null;

  renameForm: FormGroup;
  transferOwnerForm: FormGroup;

  constructor(
    public view: ResourceViewService,
    public resources: ResourceService,
    public alerts: AlertService,
    private formBuilder: FormBuilder,
    private router: Router,
  ) {
    this.renameForm = this.formBuilder.group({name: ''});
    this.transferOwnerForm = this.formBuilder.group({username: ''});
  }

  ngOnInit(): void {
    this.view.resource$.subscribe(resource => this.resource = resource);
    this.view.packages$.subscribe(packages => this.packages = packages);
  }

  submitRename({name}: {name: string}) {
    this.resources.patch(this.resource.author_id, this.resource.id, {name}).pipe(
      catchError((reason: string) => {
        this.alerts.setAlert(reason);
        return throwError(reason);
      }),
    ).subscribe(() => {
      this.router.navigate(['/u', this.resource.authors[0].username, name]);
    });
  }

  submitTransfer({username}: {username: string}) {
    if (!confirm(`Are you sure you want to transfer ownership to ` + username)) {
      return;
    }

    this.resources.transfer(this.resource.author_id, this.resource.id, username).pipe(
      catchError((reason: string) => {
        this.alerts.setAlert(reason);
        return throwError(reason);
      }),
    ).subscribe((data) => {
      this.router.navigate(['/u', data.new_username, this.resource.name]);
    });
  }

  toggleVisibility() {
    const newVisibility = (this.resource.visibility === ResourceVisibility.PRIVATE) ?
      ResourceVisibility.PUBLIC : ResourceVisibility.PRIVATE;
    this.resources.patch(this.resource.author_id, this.resource.id, {
      visibility: newVisibility,
    }).pipe(
      catchError((reason: string) => {
        this.alerts.setAlert(reason);
        return throwError(reason);
      }),
    ).subscribe(() => {
      this.resource.visibility = newVisibility;
    });
  }

}
