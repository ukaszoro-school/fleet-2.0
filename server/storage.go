package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Storage struct {
	ctx             context.Context
	db              mongo.Database
	userCollection  mongo.Collection
	routeCollection mongo.Collection
	stopCollection  mongo.Collection
}

type Driver struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     uint   `json:"age"`
}

type Stop struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Location string             `bson:"location"`
}

func storageNew(ctx context.Context, uri string) (*Storage, *mongo.Client) {

	s := Storage{}
	s.ctx = ctx
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(s.ctx, nil); err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Connected to MongoDB!")

	s.db = *client.Database("fleet")
	s.userCollection = *s.db.Collection("users")
	s.routeCollection = *s.db.Collection("routes")
	s.stopCollection = *s.db.Collection("stops")

	return &s, client
}

func (s *Storage) createStop(stop *Stop) (*mongo.InsertOneResult, error) {
	result, err := s.stopCollection.InsertOne(s.ctx, stop)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Storage) GetAll() ([]Stop, error) {
	var stops []Stop

	cursor, err := s.stopCollection.Find(s.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	for cursor.Next(s.ctx) {
		var stop Stop
		if err := cursor.Decode(&stop); err != nil {
			return nil, err
		}
		stops = append(stops, stop)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return stops, nil
}

func (s *Storage) GetByLocation(location string) (*Stop, error) {
	var stop Stop

	filter := bson.M{"location": location}

	err := s.stopCollection.FindOne(s.ctx, filter).Decode(&stop)
	if err != nil {
		return nil, err
	}

	return &stop, nil
}

func (s *Storage) DeleteByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := s.stopCollection.DeleteOne(s.ctx, filter)
	return err
}
