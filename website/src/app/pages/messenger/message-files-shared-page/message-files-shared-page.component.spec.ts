import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MessageFilesSharedPageComponent } from './message-files-shared-page.component';

describe('MessageFilesSharedPageComponent', () => {
  let component: MessageFilesSharedPageComponent;
  let fixture: ComponentFixture<MessageFilesSharedPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MessageFilesSharedPageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MessageFilesSharedPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
