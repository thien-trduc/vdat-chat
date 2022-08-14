import { TestBed } from '@angular/core/testing';

import { GenerateColorService } from './generate-color.service';

describe('GenerateColorService', () => {
  let service: GenerateColorService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GenerateColorService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
