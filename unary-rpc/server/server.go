package server

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"example.com/m/pb"
	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type User struct {
	Id   string
	Name string
	Age  int32
}

func Run() {
	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		panic(err)
	}

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(authInterceptor))
	pb.RegisterUserServer(s, NewUserService())
	reflection.Register(s)

	err = s.Serve(listen)
	if err != nil {
		panic(err)
	}
}

type UserService struct {
	pb.UnimplementedUserServer
	users map[string]*User
	mu    *sync.Mutex
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
		mu:    &sync.Mutex{},
	}
}

func (us *UserService) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	user := &User{
		Id:   req.Id,
		Name: req.Name,
		Age:  req.Age,
	}

	us.users[user.Id] = user

	return &pb.AddUserResponse{
		Id:   user.Id,
		Name: user.Name,
		Age:  user.Age,
	}, nil
}

func (us *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	us.mu.Lock()
	defer us.mu.Unlock()

	user, ok := us.users[req.Id]
	if !ok {
		return nil, errors.New("user not found")
	}

	return &pb.GetUserResponse{
		Id:   user.Id,
		Name: user.Name,
		Age:  user.Age,
	}, nil
}

func (us *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Username == "test" && req.Password == "test" {
		token, err := generateJWT(req.Username)
		if err != nil {
			return nil, err
		}

		return &pb.LoginResponse{
			Token: token,
		}, nil
	}

	return nil, errors.New("invalid credentials")
}

func generateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtKey)
}
