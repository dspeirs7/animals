export interface Animal {
  id: string;
  name: string;
  description: string;
  imageUrl: string;
  type: AnimalType;
  breed: CatBreed | ChickenBreed | DogBreed;
  vaccinations: Vaccination[];
}

export interface Vaccination {
  name: string;
  dateGiven: Date;
  dateNeeded: Date;
}

export enum AnimalType {
  Cat = 1,
  Chicken = 2,
  Dog = 3,
}

export enum CatBreed {
  RussianBlue = 1,
}

export enum ChickenBreed {
  Brahma = 10,
  BuffOrpington = 11,
}

export enum DogBreed {
  Mix = 20,
  FoxHound = 21,
}

export type AnimalName = 'Cat' | 'Chicken' | 'Dog';
