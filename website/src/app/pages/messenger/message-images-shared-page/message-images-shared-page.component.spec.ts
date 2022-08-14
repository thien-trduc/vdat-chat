import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageImagesSharedPageComponent } from './message-images-shared-page.component';

describe('MessageImagesSharedPageComponent', () => {
  let component: MessageImagesSharedPageComponent;
  let fixture: ComponentFixture<MessageImagesSharedPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageImagesSharedPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageImagesSharedPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
