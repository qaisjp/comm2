import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceUploadComponent } from './upload.component';

describe('UploadComponent', () => {
  let component: ResourceUploadComponent;
  let fixture: ComponentFixture<ResourceUploadComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceUploadComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceUploadComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
