package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"rahulchhabra.io/config"
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

// Responsible for starting the server
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

	// Start the server on port 50051
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

	// Start the server in a new goroutine (concurrency) (Serve).
	// This is so that the server can continue to run while we do other things in the main function and not block the main function.
	go func () {
		if err := grpcServer.Serve(listner); err != nil {
			log.Fatalf("Failed to serve: %s", err)
		}
	}()
	// Create a new gRPC-Gateway server (gateway).
	connection, err := grpc.DialContext(
		context.Background(),
		"localhost:50051",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	// Create a new gRPC-Gateway mux (gateway).
	gwmux := runtime.NewServeMux()

	// Register the service with the server (gateway).
	err = userproto.RegisterUserServiceHandler(context.Background(), gwmux, connection)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	// Create a new HTTP server (gateway). (Serve). (ListenAndServe)
	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}

// Responsible for creating a user.
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
	password := request.GetPassword()

	// hash the password
	hashedPassword, err := config.CreateToken(password)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Could not hash password : %s", err),
		)
	}
	// Create a new user
	newUser := model.User{
		FirstName: userdata.GetFirstname(),
		LastName:  userdata.GetLastname(),
		Email:     userdata.GetEmail(),
		Username:  userdata.GetUsername(),
		Password:  string(hashedPassword),
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
			Firstname: newUser.FirstName,
			Lastname:  newUser.LastName,
			Email:     newUser.Email,
			Username:  newUser.Username,
		},
	}, nil
}
func (*userService) AuthenticateUser(ctx context.Context, request *userproto.AuthenticateUserRequest) (response *userproto.AuthenticateUserResponse, err error) {
	// get the user details.
	email := request.Email
	password := request.Password

	// check if the user exists
	user := UserCollection.FindOne(context.Background(), model.User{Email: email})
	// check for errors
	if user.Err() != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("User with username %s not found", email),
		)
	}
	// create a user model
	var userModel model.User
	// decode the user from the database to the user struct (Decode).
	if err := user.Decode(&userModel); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Could not decode user data: %s", err),
		)
	}

	// compare user passwords(hashedpassword, inputpassword)..
	if err := config.ComparePasswords(userModel.Password, password); err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated,
			fmt.Sprintf("Password is incorrect: %s", err),
		)
	}
	// return the response
	return &userproto.AuthenticateUserResponse{
		Message: "User Authenticated Successfully",
	}, nil
}

func main() {
	// Start the server
	startServer()
}
