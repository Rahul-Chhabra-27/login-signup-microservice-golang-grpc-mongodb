package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rahulchhabra.io/model"
	userproto "rahulchhabra.io/proto/user"
)

// Create a struct that will implement the UserServiceServer interface
type userService struct {
	// This is the same as the UserServiceServer interface from the proto file (user.proto) but with an extra method called 
	// mustEmbedUnimplementedUserServiceServer() to make sure that the struct implements the UserServiceServer interface
	// This is a GoLang thing and is not required in other languages
	userproto.UnimplementedUserServiceServer
}

// Create a global variable to store the MongoDB collection
var UserCollection *mongo.Collection

func startServer() {
	// Log a message
	fmt.Println("Starting server...")
	// Create a new context
	ctx := context.TODO()

	// Connect to the MongoDB database
	db, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://chhabrarahul027:password2707@cluster.l1ycf7p.mongodb.net/"))

	// Check for errors
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	// Set the global variable to the collection
	UserCollection = db.Database("testdb").Collection("users")
	// Start the server
	listner, err := net.Listen("tcp", "localhost:50051")
	// Check for errors
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	fmt.Println("Database connected Successfully")
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Register the service with the server
	userproto.RegisterUserServiceServer(grpcServer, &userService{})
	// Check for errors
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}

// Unary Rpc -> Responsible for creating a user.
func (*userService) CreateUser(ctx context.Context, request *userproto.CreateUserRequest) (response *userproto.CreateUserResponse, err error) {
	userdata := request.GetUser()
	
	// Create a new user struct to be inserted into the database later on  (Filter).
	userfiler := model.User{
		Email: userdata.GetEmail(),
	}
	// Check if the user already exists
	user := UserCollection.FindOne(context.Background(), userfiler)
	// Check for errors
	if user.Err() == nil {
		return nil, status.Errorf(
			codes.AlreadyExists,
			fmt.Sprintf("User with email %s already exists", userdata.GetEmail()),
		)
	}
	// Create a new user
	newUser := model.User{
		Id:        primitive.NewObjectID(),
		FirstName: userdata.GetFirstName(),
		LastName:  userdata.GetLastName(),
		Email:     userdata.GetEmail(),
		Username:  userdata.GetUsername(),
		Password:  request.GetPassword(),
	}

	// Insert the user into the database
	result, err := UserCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %s", err),
		)
	}
	// Get the OID(ObjectId) of the inserted user
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID %v", err),
		)
	}
	
	return &userproto.CreateUserResponse{
		User: &userproto.User{
			Id:        oid.Hex(),
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
			Username:  newUser.Username,
		},
	}, nil
}

func main() {
	// Start the server
	startServer()
}
