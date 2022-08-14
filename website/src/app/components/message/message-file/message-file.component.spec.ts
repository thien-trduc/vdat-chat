import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageFileComponent } from './message-file.component';

describe('MessageFileComponent', () => {
  let component: MessageFileComponent;
  let fixture: ComponentFixture<MessageFileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageFileComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageFileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
