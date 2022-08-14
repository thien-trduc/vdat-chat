import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageListMemberPageComponent } from './message-list-member-page.component';

describe('MessageListMemberPageComponent', () => {
  let component: MessageListMemberPageComponent;
  let fixture: ComponentFixture<MessageListMemberPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageListMemberPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageListMemberPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
