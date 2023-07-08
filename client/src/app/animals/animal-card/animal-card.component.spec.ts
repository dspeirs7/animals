import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AnimalCardComponent } from './animal-card.component';

describe('AnimalCardComponent', () => {
  let component: AnimalCardComponent;
  let fixture: ComponentFixture<AnimalCardComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AnimalCardComponent]
    });
    fixture = TestBed.createComponent(AnimalCardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
