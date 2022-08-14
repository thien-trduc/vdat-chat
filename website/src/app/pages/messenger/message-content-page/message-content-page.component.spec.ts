import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageContentPageComponent } from './message-content-page.component';

describe('MessageContentPageComponent', () => {
  let component: MessageContentPageComponent;
  let fixture: ComponentFixture<MessageContentPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageContentPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageContentPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
