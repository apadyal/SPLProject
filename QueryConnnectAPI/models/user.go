package models

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	UserObjectId  primitive.ObjectID
	UserName      string
	UserRole      string
	UserRoleClass int32
	Email         string
	Level         int32
	Password      string
}

func init() {

}

// extractUser
func extractUser(ctx context.Context, cur *mongo.Cursor) []User {

	var users []User
	for cur.Next(ctx) {

		var u User
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			fmt.Printf("error occurred! Error is - %s", err.Error())
		}
		//fmt.Printf("%s", result["topic"])
		resultMap := result.Map()
		mapKey := resultMap["_id"].(primitive.ObjectID) //Key of the map element

		u.UserObjectId = mapKey

		u.Email = resultMap["email"].(string)

		u.UserName = resultMap["userName"].(string)
		//usrRoleCls := resultMap["userRoleClass"]
		u.UserRoleClass = resultMap["userRoleClass"].(int32)
		u.Level = resultMap["level"].(int32)
		u.UserRole = resultMap["userRole"].(string)
		u.Password = resultMap["password"].(string)

		users = append(users, u)
	}
	return users
}

// Adduser
func AddUser(user User) string {

	//var user User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Define context
	defer cancel()

	client, err := getClientConnection(ctx)
	if err != nil {
		fmt.Printf("error occurred while connecting to client! Error is - %s", err.Error())
	}
	collection := getModelCollection(client, "users")

	//uLevel, err := strconv.Atoi(user.Level)

	result, err := collection.InsertOne(ctx, bson.M{
		"userId":        rand.Intn(999999),
		"level":         2,
		"userName":      user.UserName,
		"password":      user.Password,
		"userRole":      "contributor",
		"userRoleClass": 2,
		"email":         user.Email,
		"description":   "This is a contributor role. UserRoleClass is 2 which indicates role is 'contributor', the user who conributes to the discusson.Level 2 indicates that user is a silver member"})

	if err != nil {
		//fmt.Printf("error occurred! Error is - %s", err.Error())
		return "Failure - " + err.Error()
	}

	defer client.Disconnect(ctx)
	return result.InsertedID.(primitive.ObjectID).Hex()
}

// AutheticateUser
func AutheticateUser(user User) (*User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Define context
	defer cancel()

	client, err := getClientConnection(ctx)
	if err != nil {
		fmt.Printf("error occurred while connecting to client! Error is - %s", err.Error())
	}
	collection := getModelCollection(client, "users")

	userEmail := user.Email
	filter := bson.M{"email": userEmail}

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		fmt.Printf("error occurred! Error is - %s", err.Error())
		return nil, err
	}

	users := extractUser(ctx, cur)

	if len(users) > 1 {
		return nil, errors.New("User does not exist")
	}

	if len(users) == 1 && users[0].Password != user.Password {
		return nil, errors.New("authentication failed")
	}

	if len(users) == 1 {
		//do not send password. Make it nil and return
		users[0].Password = ""
		return &users[0], nil
	}

	defer client.Disconnect(ctx)
	return nil, nil

}

// GetUserByID
func GetUserByID(id string) User {
	var users []User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //Define context
	defer cancel()

	client, err := getClientConnection(ctx)
	if err != nil {
		fmt.Printf("error occurred while connecting to client! Error is - %s", err.Error())
	}
	collection := getModelCollection(client, "users")

	//Declare filter
	userId, err := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": userId}

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		fmt.Printf("error occurred! Error is - %s", err.Error())
	}

	users = extractUser(ctx, cur)
	defer cur.Close(ctx)
	defer client.Disconnect(ctx)
	return users[0]
}
