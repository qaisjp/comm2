import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceAboutComponent } from './about.component';

describe('AboutComponent', () => {
  let component: ResourceAboutComponent;
  let fixture: ComponentFixture<ResourceAboutComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceAboutComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceAboutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
