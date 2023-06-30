import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChickenCardComponent } from './chicken-card.component';

describe('ChickenCardComponent', () => {
  let component: ChickenCardComponent;
  let fixture: ComponentFixture<ChickenCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ChickenCardComponent]
    });
    fixture = TestBed.createComponent(ChickenCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
