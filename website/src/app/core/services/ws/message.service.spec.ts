import { TestBed } from '@angular/core/testing';

import { MessageService } from './message.service';

describe('MessageV2Service', () => {
  let service: MessageService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MessageService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
