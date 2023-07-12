import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CapturedDataGoroutineComponent } from 'src/app/components/captured-data-goroutine/captured-data-goroutine.component';

describe('CapturedDataComponent', () => {
  let component: CapturedDataGoroutineComponent;
  let fixture: ComponentFixture<CapturedDataGoroutineComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [CapturedDataGoroutineComponent]
    });
    fixture = TestBed.createComponent(CapturedDataGoroutineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
