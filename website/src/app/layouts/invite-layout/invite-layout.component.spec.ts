import { ComponentFixture, TestBed } from '@angular/core/testing';

import { InviteLayoutComponent } from './invite-layout.component';

describe('InviteLayoutComponent', () => {
  let component: InviteLayoutComponent;
  let fixture: ComponentFixture<InviteLayoutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ InviteLayoutComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(InviteLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
