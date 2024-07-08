package main

import (
	"context"
	"testing"

	proto "server/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "address", "phone", "height", "married"}).
		AddRow(1, "John Doe", "123 Main St", "1234567890", 5.9, true)
	mock.ExpectQuery("SELECT id, name, address, phone, height, married FROM users WHERE id = \\$1").
		WithArgs(1).
		WillReturnRows(rows)

	s := &server{db: db}
	req := &proto.GetUserRequest{Id: 1}
	resp, err := s.GetUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, "John Doe", resp.User.Name)
	assert.Equal(t, "123 Main St", resp.User.Address)
	assert.Equal(t, int64(1234567890), resp.User.Phone)
	assert.Equal(t, float32(5.9), resp.User.Height)
	assert.True(t, resp.User.Married)
}

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "address", "phone", "height", "married"}).
		AddRow(1, "John Doe", "123 Main St", "1234567890", 5.9, true).
		AddRow(2, "Jane Doe", "456 Elm St", "9876543210", 5.7, false)
	mock.ExpectQuery("SELECT id, name, address, phone, height, married FROM users WHERE id = ANY\\(\\$1\\)").
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(rows)

	s := &server{db: db}
	req := &proto.GetUsersRequest{Ids: []int64{1, 2}}
	resp, err := s.GetUsers(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 2)
	assert.Equal(t, int64(1), resp.Users[0].Id)
	assert.Equal(t, "John Doe", resp.Users[0].Name)
	assert.Equal(t, int64(2), resp.Users[1].Id)
	assert.Equal(t, "Jane Doe", resp.Users[1].Name)
}

func TestSearchUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "address", "phone", "height", "married"}).
		AddRow(1, "John Doe", "123 Main St", "1234567890", 5.9, true)
	mock.ExpectQuery("SELECT id, name, address, phone, height, married FROM users WHERE address = \\$1 AND phone = \\$2 AND married = \\$3").
		WithArgs("123 Main St", int64(1234567890), true).
		WillReturnRows(rows)

	s := &server{db: db}
	req := &proto.SearchUsersRequest{Address: "123 Main St", Phone: 1234567890, Married: true}
	resp, err := s.SearchUsers(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Users, 1)
	assert.Equal(t, int64(1), resp.Users[0].Id)
	assert.Equal(t, "John Doe", resp.Users[0].Name)
	assert.Equal(t, "123 Main St", resp.Users[0].Address)
	assert.Equal(t, int64(1234567890), resp.Users[0].Phone)
	assert.Equal(t, float32(5.9), resp.Users[0].Height)
	assert.True(t, resp.Users[0].Married)
}
