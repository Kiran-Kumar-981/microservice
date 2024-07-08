package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	proto "server/user"

	"github.com/lib/pq"
)

// {"ID": 1, "Name": "Steve", "Address": "LA", "Phone": 1234567890, "Height": 5.8, "Married": true}
type User struct {
	ID      int64
	Name    string
	Address string
	Phone   string
	Height  float64
	Married bool
}

func NewServer() *server {
	connStr := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" host=" + os.Getenv("DB_HOST") +
		" sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &server{db: db}
}

func (s *server) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	var user proto.UserData
	err := s.db.QueryRow("SELECT id, name, address, phone, height, married FROM users WHERE id = $1", req.Id).Scan(
		&user.Id, &user.Name, &user.Address, &user.Phone, &user.Height, &user.Married,
	)
	if err != nil {
		return nil, err
	}

	return &proto.UserResponse{User: &user}, nil
}

func (s *server) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.UsersResponse, error) {
	rows, err := s.db.Query("SELECT id, name, address, phone, height, married FROM users WHERE id = ANY($1)", pq.Array(req.Ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*proto.UserData
	for rows.Next() {
		var user proto.UserData
		if err := rows.Scan(&user.Id, &user.Name, &user.Address, &user.Phone, &user.Height, &user.Married); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &proto.UsersResponse{Users: users}, nil
}

func (s *server) SearchUsers(ctx context.Context, req *proto.SearchUsersRequest) (*proto.UsersResponse, error) {
	rows, err := s.db.Query("SELECT id, name, address, phone, height, married FROM users WHERE address = $1 AND phone = $2 AND married = $3",
		req.Address, req.Phone, req.Married)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*proto.UserData
	for rows.Next() {
		var user proto.UserData
		if err := rows.Scan(&user.Id, &user.Name, &user.Address, &user.Phone, &user.Height, &user.Married); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &proto.UsersResponse{Users: users}, nil
}
