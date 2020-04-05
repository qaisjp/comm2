import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {ActivatedRoute, Router} from '@angular/router';
import {ResourcePackage, ResourceService} from '../resource.service';
import {ResourceViewService} from '../resource-view.service';
import {AlertService} from '../../alert.service';
import {FormBuilder, FormGroup} from '@angular/forms';

@Component({
  selector: 'app-upload',
  templateUrl: './upload.component.html',
  styleUrls: ['./upload.component.scss']
})
export class ResourceUploadComponent implements OnInit {
  editMode: boolean;
  pkg: ResourcePackage;
  form: FormGroup = this.formBuilder.group({
    description: '',
    draft: true,
  });

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private resources: ResourceService,
    public view: ResourceViewService,
    private alerts: AlertService,
    private formBuilder: FormBuilder,
    private cd: ChangeDetectorRef,
  ) {

  }

  ngOnInit(): void {
    this.view.uploadProgress = 0;
    this.route.data.subscribe(({ pkg }: { pkg: ResourcePackage }) => {
      this.editMode = pkg !== undefined;
      if (!this.editMode) {
        return;
      }
      this.pkg = pkg;
      this.form.setValue({
        description: pkg.description,
        draft: !pkg.published_at,
      });
      console.log(pkg);
    });
  }

  onFileChange(event) {
    if (event.target.files && event.target.files.length) {
      const file: File = event.target.files[0];

      if (this.editMode) {

      } else {
        this.view.createPackage(file).subscribe(id => {
          this.editMode = true;
          this.router.navigate(['..', 'edit', id], {
            relativeTo: this.route,
          });
        });
      }
      console.log(file);
    }
  }

}
