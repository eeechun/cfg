package services

import (
	"cfg/models"
	"context"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

type MongoDBClient struct {
	client *mongo.Client
	ctx    context.Context
}

var mongoClient = NewMongoDBClient(os.Getenv("MONGODB_URI"))

func NewMongoDBClient(uri string) *MongoDBClient {
	client, ctx, err := ConnectMongo(uri)
	if err != nil {
		log.Println("connect mongodb failed")
		log.Println(err.Error())
	}
	if PingMongo(client, ctx) != nil {
		log.Println("ping failed")
		log.Println(err.Error())
	}
	return &MongoDBClient{client, ctx}
}

func (mongoClient *MongoDBClient) FindBinary(databaseName string, collectionName string, filter bson.D) *models.Binary {
	collection := mongoClient.client.Database(databaseName).Collection(collectionName)
	var result *models.Binary
	err := collection.FindOne(mongoClient.ctx, filter).Decode(&result)
	if err != nil {
		log.Println("find data failed")
		log.Println(err.Error())
	}
	return result
}

func (mongoClient *MongoDBClient) ListBinaries(databaseName string, collectionName string, filter bson.D, projection bson.D) []models.Binary {
	collection := mongoClient.client.Database(databaseName).Collection(collectionName)
	opts := options.Find().SetProjection(projection)
	cursor, err := collection.Find(mongoClient.ctx, filter, opts)
	if err != nil {
		log.Println(err.Error())
	}
	var results []models.Binary
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	if err != nil {
		log.Println("find data failed")
		log.Println(err.Error())
	}
	return results
}

func (mongoClient *MongoDBClient) ListField(databaseName string, collectionName string, field string) []interface{} {
	collection := mongoClient.client.Database(databaseName).Collection(collectionName)
	results, err := collection.Distinct(context.TODO(), field, bson.D{})
	if err != nil {
		log.Println(err.Error())
	}
	return results
}

func (mongoClient *MongoDBClient) insertBinary(databaseName string, collectionName string, binary models.Binary) bool {
	collection := mongoClient.client.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertOne(mongoClient.ctx, binary)
	if err != nil {
		return false
	}
	return true
}

func ConnectMongo(uri string) (*mongo.Client, context.Context, error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, err
}

func CloseMongo(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func PingMongo(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	log.Println("connected successfully")
	return nil
}
