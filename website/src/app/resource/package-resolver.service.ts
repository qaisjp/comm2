import { Injectable } from '@angular/core';
import {ResourcePackage, ResourceService} from './resource.service';
import {ActivatedRouteSnapshot, Resolve, Router, RouterStateSnapshot} from '@angular/router';
import {ResourceViewService} from './resource-view.service';
import {EMPTY, Observable} from 'rxjs';
import {catchError, switchMap, take, tap} from 'rxjs/operators';
import {AlertService} from '../alert.service';

@Injectable({
  providedIn: 'root'
})
export class ResourcePackageResolverService implements Resolve<ResourcePackage> {

  constructor(
    private view: ResourceViewService,
    private resources: ResourceService,
    private router: Router,
    private alerts: AlertService,
  ) { }

  resolve(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<ResourcePackage> | Observable<never> {
    const id = Number.parseInt(route.paramMap.get('pkg_id'), 10);
    const invalidPackageID = 'Invalid package ID provided';
    if (isNaN(id)) {
      this.alerts.setAlert(invalidPackageID);
      return EMPTY;
    }

    // todo: this is kinda crappy. write a resource level resolver + figure out how to access that data from here
    // this can be (potentially) resolved by using observables less (i think i'm using it too much)
    const username = route.parent.parent.paramMap.get('username');
    const resource = route.parent.parent.paramMap.get('resource');

    return this.resources.getPackage(username, resource, id).pipe(
      tap(r => console.log('got resource', r.author_id, r.id)),
      catchError(err => {
        this.alerts.setAlert(invalidPackageID);
        this.router.navigate(['..', '..']);
        return EMPTY;
      }),
      take(1),
    );
  }
}
