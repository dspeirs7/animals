package repository

import (
	"context"

	"github.com/dspeirs7/animals/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	userColl *mongo.Collection
}

func NewUserRepository(userColl *mongo.Collection) domain.UserRepository {
	return &userRepository{userColl: userColl}
}

func (m *userRepository) GetUser(ctx context.Context, username string) (*domain.User, error) {
	filter := bson.M{"username": username}

	var user domain.User

	cursor := m.userColl.FindOne(ctx, filter)
	if err := cursor.Decode(&user); err != nil {
		return &domain.User{}, err
	}

	return &user, nil
}
