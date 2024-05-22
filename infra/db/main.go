package db

import (
	"context"
	"time"

	"github.com/gabrielmoura/WebCrawler/config"
	"github.com/gabrielmoura/WebCrawler/infra/data"
	"github.com/gabrielmoura/WebCrawler/infra/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var client *mongo.Client

func InitDB() error {
	var err error
	client, err = mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(config.Conf.MongoURI),
	)
	if err != nil {
		log.Logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
		return err
	}

	// Verifica a conexão com o banco de dados
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Logger.Fatal("Failed to ping MongoDB", zap.Error(err))
		return err
	}

	log.Logger.Info("Connected to MongoDB")
	return nil
}

func WritePage(page *data.Page) error {
	_, err := client.Database("webcrawler").Collection("pages").InsertOne(context.TODO(), page)
	if err != nil {
		log.Logger.Error("Failed to write page to MongoDB", zap.Error(err))
		return err
	}
	return nil
}

func ReadPage(url string) (*data.Page, error) {
	var page data.Page
	err := client.Database("webcrawler").Collection("pages").FindOne(context.TODO(), bson.D{{Key: "url", Value: url}}).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Logger.Info("No document found with the given URL", zap.String("url", url))
			return nil, nil
		}
		log.Logger.Error("Failed to read page from MongoDB", zap.Error(err))
		return nil, err
	}
	return &page, nil
}

func IsVisited(url string) bool {
	var page data.Page
	err := client.Database("webcrawler").Collection("pages").FindOne(context.TODO(), bson.D{{Key: "url", Value: url}}).Decode(&page)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Logger.Error("Failed to check if page is visited in MongoDB", zap.Error(err))
		return false
	}
	return true
}

func AllVisited() ([]string, error) {
	opts := options.Find().SetProjection(bson.D{{"_id", 0}, {"url", 1}})

	cursor, err := client.Database("webcrawler").Collection("pages").Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Logger.Error("Failed to retrieve all visited pages from MongoDB", zap.Error(err))
		return nil, err
	}
	defer cursor.Close(context.Background())

	var urls []string
	for cursor.Next(context.Background()) {
		var page data.Page
		if err := cursor.Decode(&page); err != nil {
			log.Logger.Error("Failed to decode page from MongoDB cursor", zap.Error(err))
			return nil, err
		}
		urls = append(urls, page.Url)
	}

	if err := cursor.Err(); err != nil {
		log.Logger.Error("Cursor error while iterating over visited pages", zap.Error(err))
		return nil, err
	}

	return urls, nil
}
