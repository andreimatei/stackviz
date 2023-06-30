import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CapturedDataComponent } from './captured-data.component';

describe('CapturedDataComponent', () => {
  let component: CapturedDataComponent;
  let fixture: ComponentFixture<CapturedDataComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [CapturedDataComponent]
    });
    fixture = TestBed.createComponent(CapturedDataComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
