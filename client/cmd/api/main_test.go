package main_test

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"client/user"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &mockUserService{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type mockUserService struct {
	user.UnimplementedUserServiceServer
}

func (s *mockUserService) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.UserResponse, error) {
	return &user.UserResponse{
		User: &user.UserData{
			Id:   req.Id,
			Name: "John Doe",
		},
	}, nil
}

func (s *mockUserService) GetUsers(ctx context.Context, req *user.GetUsersRequest) (*user.UsersResponse, error) {
	users := []*user.UserData{
		{Id: 1, Name: "John Doe"},
		{Id: 2, Name: "Jane Doe"},
	}
	return &user.UsersResponse{Users: users}, nil
}

func (s *mockUserService) SearchUsers(ctx context.Context, req *user.SearchUsersRequest) (*user.UsersResponse, error) {
	users := []*user.UserData{
		{Id: 1, Name: "John Doe", Address: req.Address, Phone: req.Phone, Married: req.Married},
	}
	return &user.UsersResponse{Users: users}, nil
}

func TestGetUser(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufconn: %v", err)
	}
	defer conn.Close()
	client := user.NewUserServiceClient(conn)

	handler := func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/user/")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		req := &user.GetUserRequest{Id: id}
		res, err := client.GetUser(context.Background(), req)
		if err != nil {
			http.Error(w, "Failed to get user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.User)
	}

	req := httptest.NewRequest("GET", "/user/1", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}

	var got user.UserData
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	want := &user.UserData{
		Id:   1,
		Name: "John Doe",
	}

	if diff := cmp.Diff(want, &got, protocmp.Transform()); diff != "" {
		t.Errorf("Unexpected users data (-want +got):\n%s", diff)
	}
}

func TestGetUsers(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufconn: %v", err)
	}
	defer conn.Close()
	client := user.NewUserServiceClient(conn)

	handler := func(w http.ResponseWriter, r *http.Request) {
		ids := r.URL.Query()["ids"]
		var intIds []int64
		for _, id := range ids {
			intId, err := strconv.ParseInt(id, 10, 64)
			if err != nil {
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}
			intIds = append(intIds, intId)
		}
		req := &user.GetUsersRequest{Ids: intIds}
		res, err := client.GetUsers(context.Background(), req)
		if err != nil {
			http.Error(w, "Failed to get users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Users)
	}

	req := httptest.NewRequest("GET", "/users?ids=1&ids=2", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}

	var got []*user.UserData
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	want := []*user.UserData{
		{Id: 1, Name: "John Doe"},
		{Id: 2, Name: "Jane Doe"},
	}

	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("Unexpected users data (-want +got):\n%s", diff)
	}
}

func TestSearchUsers(t *testing.T) {
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufconn: %v", err)
	}
	defer conn.Close()
	client := user.NewUserServiceClient(conn)

	handler := func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		phoneStr := r.URL.Query().Get("phone")
		marriedStr := r.URL.Query().Get("married")
		phone, err := strconv.ParseInt(phoneStr, 10, 64)
		if err != nil {
			phone = 0 // Handle missing or invalid phone parameter
		}
		married, err := strconv.ParseBool(marriedStr)
		if err != nil {
			married = false // Handle missing or invalid married parameter
		}
		req := &user.SearchUsersRequest{
			Address: address,
			Phone:   phone,
			Married: married,
		}
		res, err := client.SearchUsers(context.Background(), req)
		if err != nil {
			http.Error(w, "Failed to search users: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Users)
	}

	req := httptest.NewRequest("GET", "/search?address=123%20Main%20St&phone=1234567890&married=true", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", resp.Status)
	}

	var got []*user.UserData
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	want := []*user.UserData{
		{Id: 1, Name: "John Doe", Address: "123 Main St", Phone: 1234567890, Married: true},
	}

	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("Unexpected users data (-want +got):\n%s", diff)
	}
}
