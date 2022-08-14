import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageEditInfoPageComponent } from './message-edit-info-page.component';

describe('MessageEditInfoPageComponent', () => {
  let component: MessageEditInfoPageComponent;
  let fixture: ComponentFixture<MessageEditInfoPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageEditInfoPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageEditInfoPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
