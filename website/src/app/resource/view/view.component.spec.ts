import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceViewComponent } from './view.component';

describe('ResourcePageComponent', () => {
  let component: ResourceViewComponent;
  let fixture: ComponentFixture<ResourceViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceViewComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
