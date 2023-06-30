import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChickenComponent } from './chicken.component';

describe('ChickenComponent', () => {
  let component: ChickenComponent;
  let fixture: ComponentFixture<ChickenComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ChickenComponent]
    });
    fixture = TestBed.createComponent(ChickenComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
