package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Animal struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	ImageUrl     string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Type         AnimalType         `bson:"type,omitempty" json:"type,omitempty"`
	Breed        AnimalBreed        `bson:"breed,omitempty" json:"breed,omitempty"`
	Vaccinations []Vaccination      `bson:"vaccinations,omitempty" json:"vaccinations,omitempty"`
}

type AnimalType int

const (
	CatType     AnimalType = 1
	ChickenType AnimalType = 2
	DogType     AnimalType = 3
)

type AnimalBreed int

type Vaccination struct {
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	DateGiven  primitive.DateTime `bson:"dateGiven,omitempty" json:"dateGiven,omitempty"`
	DateNeeded primitive.DateTime `bson:"dateNeeded,omitempty" json:"dateNeeded,omitEmpty"`
}

type AnimalRepository interface {
	GetAllCats(ctx context.Context) ([]*Animal, error)
	GetAllChickens(ctx context.Context) ([]*Animal, error)
	GetAllDogs(ctx context.Context) ([]*Animal, error)
	GetById(ctx context.Context, id string) (*Animal, error)
	Insert(ctx context.Context, insert Animal) (*Animal, error)
	Update(ctx context.Context, id string, update Animal) error
	AddVaccinations(ctx context.Context, id string, vaccinations []Vaccination) error
	DeleteVaccination(ctx context.Context, id string, vaccination Vaccination) error
	Delete(ctx context.Context, id string) error
	UpdateImageUrl(ctx context.Context, id string, url string) error
}
