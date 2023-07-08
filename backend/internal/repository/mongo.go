package repository

import (
	"context"
	"os"

	"github.com/dspeirs7/animals/internal/domain"
	"github.com/dspeirs7/animals/internal/log"
	secrets "github.com/ijustfool/docker-secrets"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func GetDB(ctx context.Context) *mongo.Client {
	logger := log.NewLogger("mongo")
	defer logger.Sync()

	var adminPassword, uri string

	dockerSecrets, _ := secrets.NewDockerSecrets("")

	adminPassword, err := dockerSecrets.Get("admin_password")
	if err != nil {
		adminPassword = os.Getenv("ADMIN_PASSWORD")
	}

	uri, err = dockerSecrets.Get("db_string")
	if err != nil {
		uri = os.Getenv("MONGODB_URI")
	}

	if uri == "" {
		logger.Fatal("you must set MONGODB_URI")
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		logger.Panic("error connecting", zap.Error(err))
	}

	createAdminUser(client.Database("animals").Collection("users"), adminPassword, logger)

	return client
}

func createAdminUser(userColl *mongo.Collection, adminPassword string, logger *zap.Logger) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cursor, err := userColl.Find(ctx, bson.D{})
	if err != nil {
		logger.Fatal("error getting user collection", zap.Error(err))
	}

	var users []domain.User

	err = cursor.All(context.Background(), &users)
	if err != nil {
		cancel()
		logger.Fatal("error getting users", zap.Error(err))
	}

	if len(users) == 1 {
		return
	} else if len(users) > 1 {
		userColl.DeleteMany(ctx, nil)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		cancel()
		logger.Fatal("not able to hash password", zap.Error(err))
	}

	user := domain.User{Username: "admin", Password: string(hashedPassword)}

	userColl.InsertOne(ctx, &user)
}
