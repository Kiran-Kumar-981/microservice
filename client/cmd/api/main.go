package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"client/user"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()

	router.GET("/user/:id", func(c *gin.Context) {
		id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
		req := &user.GetUserRequest{Id: id}
		res, err := client.GetUser(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.User)
	})

	router.GET("/users", func(c *gin.Context) {
		ids := c.QueryArray("ids")
		var intIds []int64
		for _, id := range ids {
			intId, _ := strconv.ParseInt(id, 10, 64)
			intIds = append(intIds, intId)
		}
		req := &user.GetUsersRequest{Ids: intIds}
		res, err := client.GetUsers(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Users)
	})

	router.GET("/search", func(c *gin.Context) {
		address := c.Query("address")
		phone, _ := strconv.ParseInt(c.Query("phone"), 10, 64)
		married, _ := strconv.ParseBool(c.Query("married"))
		req := &user.SearchUsersRequest{
			Address: address,
			Phone:   phone,
			Married: married,
		}
		res, err := client.SearchUsers(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res.Users)
	})

	router.Run(":8081")
}
