import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageListPendingMembersPageComponent } from './message-list-pending-members-page.component';

describe('MessageListPendingMembersPageComponent', () => {
  let component: MessageListPendingMembersPageComponent;
  let fixture: ComponentFixture<MessageListPendingMembersPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageListPendingMembersPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageListPendingMembersPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
