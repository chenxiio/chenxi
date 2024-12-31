package runtime

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/tj/assert"
)

func TestUsers_RemoveUser(t *testing.T) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	users := &Users{db: db}
	// Insert a user into the database
	_, err = db.Exec("INSERT INTO users (username, fullname, password, groups) VALUES ('testuser', 'Test User', 'password', 'group')")
	if err != nil {
		fmt.Println(err)
	}
	// Remove the user from the database
	err = users.RemoveUser(context.TODO(), "testuser")
	assert.NoError(t, err)
	// Check if the user is removed from the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'testuser'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
func TestUsers_SetUser(t *testing.T) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	users := &Users{db: db}
	// Set a user in the database
	err = users.SetUser(context.TODO(), "testuser", "Test User", "password", "group")
	assert.NoError(t, err)
	// Check if the user is set in the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'testuser'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
func TestUsers_GetUsers(t *testing.T) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	users := &Users{db: db}
	// Insert some users into the database
	_, err = db.Exec("INSERT INTO users (username, fullname, password, groups) VALUES ('user1', 'User 1', 'password1', 'group1')")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("INSERT INTO users (username, fullname, password, groups) VALUES ('user2', 'User 2', 'password2', 'group2')")
	if err != nil {
		fmt.Println(err)
	}
	// Get all users from the database
	data, err := users.GetUsers(context.TODO(), nil)
	assert.NoError(t, err)
	fmt.Println(len(data))
	// Get a specific user from the database
	data, err = users.GetUsers(context.TODO(), &User{Username: "user1"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(data))
	assert.Equal(t, "user1", data[0].Username)
	assert.Equal(t, "User 1", data[0].Fullname)
	assert.Equal(t, "group1", data[0].Groups)
}
func TestUsers_Init(t *testing.T) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	users := &Users{db: db}
	// Initialize the users table in the database
	err = users.Init()
	assert.NoError(t, err)
	// Check if the users table is created in the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}
func TestUsers_SetDefault(t *testing.T) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	users := &Users{db: db}
	// Set the default user in the database
	err = users.SetDefault(context.TODO())
	assert.NoError(t, err)
	// Check if the default user is set in the database
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = 'admin'").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}
