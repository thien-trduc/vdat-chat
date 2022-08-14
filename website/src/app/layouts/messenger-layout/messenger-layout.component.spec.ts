import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessengerLayoutComponent } from './messenger-layout.component';

describe('MessengerLayoutComponent', () => {
  let component: MessengerLayoutComponent;
  let fixture: ComponentFixture<MessengerLayoutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessengerLayoutComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessengerLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
