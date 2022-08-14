import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageQrCodePageComponent } from './message-qr-code-page.component';

describe('MessageQrCodePageComponent', () => {
  let component: MessageQrCodePageComponent;
  let fixture: ComponentFixture<MessageQrCodePageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageQrCodePageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageQrCodePageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
