import { TestBed } from '@angular/core/testing';

import { ResourcePackageResolverService } from './package-resolver.service';

describe('ResourcePackageResolverService', () => {
  let service: ResourcePackageResolverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ResourcePackageResolverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
