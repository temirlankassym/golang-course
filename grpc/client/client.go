package client

import (
	"context"
	pb "course/grpc/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"time"
)

type Client struct {
	conn pb.UserServiceClient
	ctx  context.Context
}

func NewClient() (*Client, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewUserServiceClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	return &Client{conn: c, ctx: ctx}, nil
}

func (c *Client) AddUser(name string, email string) (*pb.UserId, error) {
	return c.conn.AddUser(c.ctx, &pb.User{Name: name, Email: email})
}

func (c *Client) GetUser(id int32) (*pb.User, error) {
	return c.conn.GetUser(c.ctx, &pb.UserId{Id: id})
}

func (c *Client) ListUsers() (*pb.Users, error) {
	return c.conn.ListUsers(c.ctx, &emptypb.Empty{})
}

// change package to main in order to just run client
// change package to client in order to run tests
func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Printf("Adding user:")

	//Adding user
	id, err := c.AddUser(ctx, &pb.User{Name: "Temirlan", Email: "@gmail.com"})
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}
	log.Printf("Id of User: %d", id.GetId())

	log.Printf("Getting user:")

	//Getting user
	user, err := c.GetUser(ctx, &pb.UserId{Id: id.Id})
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User: %d,%s,%s", user.GetId(), user.GetName(), user.GetEmail())

	log.Printf("List of users")

	//List of users
	users, err := c.ListUsers(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not get users: %v", err)
	}

	for _, u := range users.GetUsers() {
		log.Printf("User: %d,%s,%s", u.GetId(), u.GetName(), u.GetEmail())
	}
}
