import { TestBed } from '@angular/core/testing';

import { ChickenService } from './chicken.service';

describe('ChickenService', () => {
  let service: ChickenService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ChickenService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
