import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageInfoPageComponent } from './message-info-page.component';

describe('MessageInfoPageComponent', () => {
  let component: MessageInfoPageComponent;
  let fixture: ComponentFixture<MessageInfoPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageInfoPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageInfoPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
