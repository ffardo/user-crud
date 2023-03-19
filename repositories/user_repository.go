package repositories

import (
	"context"
	"time"

	"github.com/ffardo/user-crud/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DB_TIMEOUT = 5

type UserRepository struct {
	Client *mongo.Client
}

func (u UserRepository) CreateUser(user models.User) (models.User, error) {
	collection := u.Client.Database("user_service").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT*time.Second)
	defer cancel()

	now := time.Now()
	user.UUID = uuid.New()
	user.Created = now
	user.Updated = now

	doc, err := bson.Marshal(user)

	_, err = collection.InsertOne(ctx, doc)

	return user, err
}

func (u UserRepository) GetUserByUUID(user_uuid uuid.UUID) (models.User, error) {
	filter := bson.D{{Key: "uuid", Value: user_uuid}}
	return u.getUserByFilter(filter)

}

func (u UserRepository) getUserByFilter(filter primitive.D) (models.User, error) {
	collection := u.Client.Database("user_service").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT*time.Second)

	defer cancel()

	var user models.User

	err := collection.FindOne(ctx, filter).Decode(&user)

	return user, err
}

func (u UserRepository) UserExistsWithEmail(email string) (bool, error) {
	filter := bson.D{{Key: "email", Value: email}}
	return u.userExistsByFilter(filter)
}

func (u UserRepository) UserExistsWithEmailAndNotUuid(email string, uuid uuid.UUID) (bool, error) {

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{{Key: "email", Value: email}},
				bson.D{{Key: "uuid", Value: bson.D{{Key: "$ne", Value: uuid}}}},
			},
		},
	}

	return u.userExistsByFilter(filter)
}

func (u UserRepository) userExistsByFilter(filter primitive.D) (bool, error) {

	collection := u.Client.Database("user_service").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT*time.Second)

	defer cancel()

	total, err := collection.CountDocuments(ctx, filter)

	return total > 0, err
}

func (u UserRepository) UpdateUser(user models.User) (models.User, error) {
	collection := u.Client.Database("user_service").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT*time.Second)

	defer cancel()
	filter := bson.D{{Key: "uuid", Value: user.UUID}}

	user.Updated = time.Now()

	doc, err := bson.Marshal(user)

	_, err = collection.ReplaceOne(ctx, filter, doc)
	return user, err
}

func (u UserRepository) DeleteUser(user_uuid uuid.UUID) error {
	collection := u.Client.Database("user_service").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT*time.Second)

	defer cancel()
	filter := bson.D{{Key: "uuid", Value: user_uuid}}

	_, err := collection.DeleteOne(ctx, filter)
	return err
}

func (u UserRepository) Init() error {
	collection := u.Client.Database("user_service").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT*time.Second)

	defer cancel()

	index := mongo.IndexModel{
		Keys: bson.M{
			"email": 1, // index in ascending order
		}, Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, index)

	return err

}
