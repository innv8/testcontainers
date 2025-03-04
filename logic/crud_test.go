package logic

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sql.DB

func setUpTestContainer() (*sql.DB, func(), error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpass",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL").WithOccurrence(1),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, err
	}

	host, err := mysqlC.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		return nil, nil, err
	}

	uri := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb?parseTime=true", host, port.Port())

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, nil, err
	}

	// create users table
	var query = `
	CREATE TABLE IF NOT EXISTS user (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(30) NOT NULL,
		age INT NOT NULL
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		mysqlC.Terminate(ctx)
	}
	return db, cleanup, nil
}

func TestMain(m *testing.M) {
	var cleanup func()
	var err error

	testDB, cleanup, err = setUpTestContainer()
	if err != nil {
		log.Fatalf("Failed to set up test container; %v", err)
	}
	defer cleanup()

	m.Run()
}

func TestCreateUser(t *testing.T) {
	userID, err := CreateUser(testDB, "Alice", 45)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := GetUser(testDB, userID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if user.Name != "Alice" || user.Age != 45 {
		t.Errorf("Expected Alice, got %+v", user)
	}
}

func TestGetUser(t *testing.T) {
	userID, err := CreateUser(testDB, "Bob", 12)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err := GetUser(testDB, userID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if user.Name != "Bob" || user.Age != 12 {
		t.Errorf("Expected Bob, got %+v", user)
	}
}

func TestUpdateUser(t *testing.T) {
	userID, err := CreateUser(testDB, "Charlie", 32)
	if err != nil {
		t.Fatalf("Failed to create user %v", err)
	}

	err = UpdateUser(testDB, userID, "Charlie Brown", 20)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	user, err := GetUser(testDB, userID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if user.Name != "Charlie Brown" || user.Age != 20 {
		t.Errorf("Update failed, got %+v", user)
	}
}

func TestDeleteUser(t *testing.T) {
	userID, err := CreateUser(testDB, "David", 45)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	err = DeleteUser(testDB, userID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	user, err := GetUser(testDB, userID)
	if err == nil {
		t.Fatalf("Expected user to be deleted but got %+v", user)
	}
}
