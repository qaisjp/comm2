import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminBansComponent } from './bans.component';

describe('BansComponent', () => {
  let component: AdminBansComponent;
  let fixture: ComponentFixture<AdminBansComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AdminBansComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AdminBansComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
