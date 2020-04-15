import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VoiceComponent } from './voice.component';
import { ChangeDetectorRef} from '@angular/core';
import { DomSanitizer } from '@angular/platform-browser';

describe('VoiceComponent', () => {
  let component: VoiceComponent;
  let fixture: ComponentFixture<VoiceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VoiceComponent ],
      providers: [
        { provide: ChangeDetectorRef},
        { provide: DomSanitizer }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VoiceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
