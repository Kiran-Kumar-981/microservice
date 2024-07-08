package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"client/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("grpc_server:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := user.NewUserServiceClient(conn)

	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
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
	})

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
