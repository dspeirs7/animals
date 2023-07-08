package repository

import (
	"context"

	"github.com/dspeirs7/animals/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoAnimalRepository struct {
	animalColl *mongo.Collection
}

func NewAnimalRepository(animalColl *mongo.Collection) domain.AnimalRepository {

	return &mongoAnimalRepository{
		animalColl: animalColl,
	}
}

func (m *mongoAnimalRepository) GetAllCats(ctx context.Context) ([]*domain.Animal, error) {
	filter := bson.D{{Key: "type", Value: 1}}
	cursor, err := m.animalColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*domain.Animal

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mongoAnimalRepository) GetAllChickens(ctx context.Context) ([]*domain.Animal, error) {
	filter := bson.D{{Key: "type", Value: 2}}
	cursor, err := m.animalColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*domain.Animal

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mongoAnimalRepository) GetAllDogs(ctx context.Context) ([]*domain.Animal, error) {
	filter := bson.D{{Key: "type", Value: 3}}
	cursor, err := m.animalColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []*domain.Animal

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mongoAnimalRepository) GetById(ctx context.Context, id string) (*domain.Animal, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &domain.Animal{}, err
	}

	var result domain.Animal

	cursor := m.animalColl.FindOne(ctx, bson.M{"_id": objectId})
	cursor.Decode(&result)

	return &result, nil
}

func (m *mongoAnimalRepository) Insert(ctx context.Context, animal domain.Animal) (*domain.Animal, error) {
	result, err := m.animalColl.InsertOne(ctx, animal)
	if err != nil {
		return &animal, err
	}

	animal.Id = result.InsertedID.(primitive.ObjectID)

	return &animal, nil
}

func (m *mongoAnimalRepository) Update(ctx context.Context, id string, animal domain.Animal) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	if _, err := m.animalColl.ReplaceOne(ctx, filter, animal); err != nil {
		return err
	}

	return nil
}

func (m *mongoAnimalRepository) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	if _, err := m.animalColl.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (m *mongoAnimalRepository) AddVaccinations(ctx context.Context, id string, vaccinations []domain.Vaccination) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	change := bson.M{"$push": bson.M{"vaccinations": bson.M{"$each": vaccinations}}}

	if _, err := m.animalColl.UpdateByID(ctx, objectId, change); err != nil {
		return err
	}

	return nil
}

func (m *mongoAnimalRepository) DeleteVaccination(ctx context.Context, id string, vaccination domain.Vaccination) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	change := bson.M{"$pull": bson.M{"vaccinations": bson.M{"name": vaccination.Name, "dateGiven": vaccination.DateGiven, "dateNeeded": vaccination.DateNeeded}}}

	if _, err := m.animalColl.UpdateByID(ctx, objectId, change); err != nil {
		return err
	}

	return nil
}

func (m *mongoAnimalRepository) UpdateImageUrl(ctx context.Context, id string, imageUrl string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	change := bson.M{"$set": bson.M{"imageUrl": imageUrl}}

	if _, err := m.animalColl.UpdateByID(ctx, objectId, change); err != nil {
		return err
	}

	return nil
}
