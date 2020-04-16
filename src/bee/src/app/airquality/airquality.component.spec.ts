import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AirqualityComponent } from './airquality.component';
import { FormBuilder } from '@angular/forms';
import { HttpClient, HttpHandler } from '@angular/common/http';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';

describe('AirqualityComponent', () => {
  let component: AirqualityComponent;
  let fixture: ComponentFixture<AirqualityComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AirqualityComponent ],
      providers: [
        {provide: FormBuilder},
        {provide: HttpClient},
        {provide: HttpHandler}
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AirqualityComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
