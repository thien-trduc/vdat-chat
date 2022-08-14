import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageVideoComponent } from './message-video.component';

describe('MessageVideoComponent', () => {
  let component: MessageVideoComponent;
  let fixture: ComponentFixture<MessageVideoComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageVideoComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageVideoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
