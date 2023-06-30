package repository

import (
	"context"

	"github.com/dspeirs7/animals/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoChickenRepository struct {
	chickenColl *mongo.Collection
	userColl    *mongo.Collection
}

func NewChickenRepository(chickenColl, userColl *mongo.Collection) domain.ChickenRepository {

	return &mongoChickenRepository{
		chickenColl: chickenColl,
		userColl:    userColl,
	}
}

func (m *mongoChickenRepository) GetUser(ctx context.Context, username string) (*domain.User, error) {
	filter := bson.M{"username": username}

	var user domain.User

	cursor := m.userColl.FindOne(ctx, filter)
	if err := cursor.Decode(&user); err != nil {
		return &domain.User{}, err
	}

	return &user, nil
}

func (m *mongoChickenRepository) GetAll(ctx context.Context) ([]*domain.Chicken, error) {
	cursor, err := m.chickenColl.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var results []*domain.Chicken

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (m *mongoChickenRepository) GetById(ctx context.Context, id string) (*domain.Chicken, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &domain.Chicken{}, err
	}

	var result domain.Chicken

	cursor := m.chickenColl.FindOne(ctx, bson.M{"_id": objectId})
	cursor.Decode(&result)

	return &result, nil
}

func (m *mongoChickenRepository) Insert(ctx context.Context, chicken domain.Chicken) (*domain.Chicken, error) {
	result, err := m.chickenColl.InsertOne(ctx, chicken)
	if err != nil {
		return &chicken, err
	}

	chicken.Id = result.InsertedID.(primitive.ObjectID)

	return &chicken, nil
}

func (m *mongoChickenRepository) Update(ctx context.Context, id string, chicken domain.Chicken) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	if _, err := m.chickenColl.ReplaceOne(ctx, filter, chicken); err != nil {
		return err
	}

	return nil
}

func (m *mongoChickenRepository) Delete(ctx context.Context, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	if _, err := m.chickenColl.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (m *mongoChickenRepository) AddVaccinations(ctx context.Context, id string, vaccinations []domain.Vaccination) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	change := bson.M{"$push": bson.M{"vaccinations": bson.M{"$each": vaccinations}}}

	if _, err := m.chickenColl.UpdateByID(ctx, objectId, change); err != nil {
		return err
	}

	return nil
}

func (m *mongoChickenRepository) DeleteVaccination(ctx context.Context, id string, vaccination domain.Vaccination) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	change := bson.M{"$pull": bson.M{"vaccinations": bson.M{"name": vaccination.Name, "dateGiven": vaccination.DateGiven, "dateNeeded": vaccination.DateNeeded}}}

	if _, err := m.chickenColl.UpdateByID(ctx, objectId, change); err != nil {
		return err
	}

	return nil
}

func (m *mongoChickenRepository) UpdateImageUrl(ctx context.Context, id string, imageUrl string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	change := bson.M{"$set": bson.D{{Key: "imageUrl", Value: imageUrl}}}

	if _, err := m.chickenColl.UpdateByID(ctx, objectId, change); err != nil {
		return err
	}

	return nil
}
