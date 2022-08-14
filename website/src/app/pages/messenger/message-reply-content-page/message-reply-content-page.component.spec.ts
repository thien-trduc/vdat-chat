import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageReplyContentPageComponent } from './message-reply-content-page.component';

describe('MessageReplyContentPageComponent', () => {
  let component: MessageReplyContentPageComponent;
  let fixture: ComponentFixture<MessageReplyContentPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageReplyContentPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageReplyContentPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
