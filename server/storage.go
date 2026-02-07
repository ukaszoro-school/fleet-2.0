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

type Stop struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Location string             `bson:"location"`
}

type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	name string             `bson:"name"`
}

type Route struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Times map[string]string  `bson:"times"`
}

func storageNew(ctx *context.Context, uri string) (*Storage, *mongo.Client) {

	s := Storage{}
	s.ctx = *ctx
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

func (s *Storage) getAllStops() ([]Stop, error) {
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

func (s *Storage) getStopByID(id primitive.ObjectID) (*Stop, error) {
	var stop Stop

	filter := bson.M{"_id": id}

	err := s.stopCollection.FindOne(s.ctx, filter).Decode(&stop)
	if err != nil {
		return nil, err
	}

	return &stop, nil
}

func (s *Storage) getStopByLocation(location string) (*Stop, error) {
	var stop Stop

	filter := bson.M{"location": location}

	err := s.stopCollection.FindOne(s.ctx, filter).Decode(&stop)
	if err != nil {
		return nil, err
	}

	return &stop, nil
}

func (s *Storage) deleteStopByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := s.stopCollection.DeleteOne(s.ctx, filter)
	return err
}

func (s *Storage) createUser(user *User) (*mongo.InsertOneResult, error) {
	result, err := s.stopCollection.InsertOne(s.ctx, user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Storage) getAllUsers() ([]User, error) {
	var users []User

	cursor, err := s.userCollection.Find(s.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	for cursor.Next(s.ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Storage) getUserByID(id primitive.ObjectID) (*User, error) {
	var user User

	filter := bson.M{"_id": id}

	err := s.userCollection.FindOne(s.ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) deleteUserByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := s.userCollection.DeleteOne(s.ctx, filter)
	return err
}

func (s *Storage) createRoute(route *Route) (*mongo.InsertOneResult, error) {
	result, err := s.routeCollection.InsertOne(s.ctx, route)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Storage) getAllRoutes() ([]Route, error) {
	var routes []Route

	cursor, err := s.routeCollection.Find(s.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(s.ctx)

	for cursor.Next(s.ctx) {
		var route Route
		if err := cursor.Decode(&route); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (s *Storage) deleteRouteByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := s.routeCollection.DeleteOne(s.ctx, filter)
	return err
}
