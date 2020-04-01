import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceVersionsComponent } from './versions.component';

describe('VersionsComponent', () => {
  let component: ResourceVersionsComponent;
  let fixture: ComponentFixture<ResourceVersionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceVersionsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceVersionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
