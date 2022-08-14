import { TestBed } from '@angular/core/testing';

import { WorkerUpdateService } from './worker-update.service';

describe('WorkerLogUpdateService', () => {
  let service: WorkerUpdateService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WorkerUpdateService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
