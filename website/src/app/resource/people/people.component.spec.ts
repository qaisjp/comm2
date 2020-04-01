import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ResourcePeopleComponent } from './people.component';

describe('PeopleComponent', () => {
  let component: ResourcePeopleComponent;
  let fixture: ComponentFixture<ResourcePeopleComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ResourcePeopleComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ResourcePeopleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
