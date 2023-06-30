package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chicken struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name,omitempty" json:"name,omitempty"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	ImageUrl     string             `bson:"imageUrl,omitempty" json:"imageUrl,omitempty"`
	Type         ChickenType        `bson:"type,omitempty" json:"type,omitempty"`
	Vaccinations []Vaccination      `bson:"vaccinations,omitempty" json:"vaccinations,omitempty"`
}

type ChickenType int

const (
	Brahma        ChickenType = 1
	BuffOrpington ChickenType = 2
)

type Vaccination struct {
	Name       string             `bson:"name,omitempty" json:"name,omitempty"`
	DateGiven  primitive.DateTime `bson:"dateGiven,omitempty" json:"dateGiven,omitempty"`
	DateNeeded primitive.DateTime `bson:"dateNeeded,omitempty" json:"dateNeeded,omitEmpty"`
}

type ChickenRepository interface {
	GetUser(ctx context.Context, username string) (*User, error)
	GetAll(ctx context.Context) ([]*Chicken, error)
	GetById(ctx context.Context, id string) (*Chicken, error)
	Insert(ctx context.Context, chicken Chicken) (*Chicken, error)
	Update(ctx context.Context, id string, chicken Chicken) error
	AddVaccinations(ctx context.Context, id string, vaccinations []Vaccination) error
	DeleteVaccination(ctx context.Context, id string, vaccination Vaccination) error
	Delete(ctx context.Context, id string) error
	UpdateImageUrl(ctx context.Context, id string, url string) error
}
