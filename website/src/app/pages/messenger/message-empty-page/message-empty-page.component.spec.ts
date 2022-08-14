import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageEmptyPageComponent } from './message-empty-page.component';

describe('MessageEmptyPageComponent', () => {
  let component: MessageEmptyPageComponent;
  let fixture: ComponentFixture<MessageEmptyPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageEmptyPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageEmptyPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
