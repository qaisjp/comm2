import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourceReviewsComponent } from './reviews.component';

describe('ReviewsComponent', () => {
  let component: ResourceReviewsComponent;
  let fixture: ComponentFixture<ResourceReviewsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourceReviewsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourceReviewsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
