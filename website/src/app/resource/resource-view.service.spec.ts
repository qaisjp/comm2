import { TestBed } from '@angular/core/testing';

import { ResourceViewService } from './resource-view.service';

describe('ResourceViewService', () => {
  let service: ResourceViewService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ResourceViewService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
