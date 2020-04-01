import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceLayoutComponent } from './layout.component';

describe('ResourcePageComponent', () => {
  let component: ResourceLayoutComponent;
  let fixture: ComponentFixture<ResourceLayoutComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceLayoutComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
