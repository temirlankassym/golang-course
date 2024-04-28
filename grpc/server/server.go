package main

import (
	"context"
	repo "course/grpc/server/repository"
	pb "course/grpc/user"
	"fmt"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

var db *pgx.Conn

// adding user to database
func (s *server) AddUser(ctx context.Context, in *pb.User) (*pb.UserId, error) {
	id, err := repo.AddUser(ctx, db, in.Name, in.Email)
	if err != nil {
		return nil, fmt.Errorf("cannot add user")
	}

	return &pb.UserId{Id: id}, nil
}

// getting user from database
func (s *server) GetUser(ctx context.Context, in *pb.UserId) (*pb.User, error) {
	user, err := repo.GetUser(ctx, db, in.Id)
	if err != nil {
		return nil, fmt.Errorf("cannot get user")
	}

	return &pb.User{Id: user.Id, Name: user.Name, Email: user.Email}, nil
}

// getting list of users from database
func (s *server) ListUsers(ctx context.Context, in *emptypb.Empty) (*pb.Users, error) {
	users, err := repo.ListUsers(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("cannot get users: %w", err)
	}

	list := []*pb.User{}

	for _, user := range users {
		list = append(list, &pb.User{Id: user.Id, Name: user.Name, Email: user.Email})
	}

	return &pb.Users{Users: list}, nil
}

func main() {
	ctx := context.Background()
	var err error

	//connecting to database
	db, err = repo.Connect(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close(ctx)

	//listening on port 50051 default /net package
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
