package redisutil

import (
	"testing"

	"github.com/alicebob/miniredis"
)

var valid_users = map[string]string{
	"bobman12": `{"username": "bobman12", "email": "bob@bobmail.com"}`,
	"jcdenton": `{"username": "jcdenton", "email": "jc@unatco.org"}`,
	"herpderp": `{"username": "herpderp", "email": "herp@derp.io"}`,
}

var invalid_users = map[string]string{
	"missingusername": `{"email": "missing@username.com"}`,
	"missingemail":    `{"username": "missingemail"}`,
}

func Test_GetUser(t *testing.T) {
	// Set up minikube for testing, fail if not working
	miniredis_socket, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniredis_socket.Close()

	// Set users in to test on
	for key, val := range valid_users {
		miniredis_socket.HSet("users", key, val)
	}
	// Start client
	redis_client, _ := NewRedisHashConn(miniredis_socket.Addr(), "", 0)

	for key, expected_val := range valid_users {
		actual_val, err := redis_client.GetUser(key)

		if err != nil {
			t.Logf("err: %s", err)
		}

		if expected_val != actual_val {
			t.Logf("Expected:\t %s \nGot:\t %s\n", expected_val, actual_val)
		}
	}
}

func Test_CreateUser(t *testing.T) {
	// Set up minikube for testing, fail if not working
	miniredis_socket, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniredis_socket.Close()

	redis_client, _ := NewRedisHashConn(miniredis_socket.Addr(), "", 0)

	for username, userdata := range valid_users {
		redis_client.CreateUser(username, userdata)
	}

	for key, expected_val := range valid_users {
		actual_val, err := redis_client.GetUser(key)

		if err != nil {
			t.Logf("err: %s", err)
		}

		if expected_val != actual_val {
			t.Logf("Expected:\t %s \nGot:\t %s\n", expected_val, actual_val)
		}
	}
}

/*
func Test_DeleteUser(t *testing.T) {
	// Set up minikube for testing, fail if not working
	miniredis_socket, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer miniredis_socket.Close()

	redis_client := NewRedisHashConn(miniredis_socket.Addr(), "", 0)

	for _, val := range valid_users {
		redis_client.CreateUser(val)
	}

	for key, _ := range valid_users {

		, err := redis_client.DeleteUser(key)

		if err != nil {
			t.Logf("err: %s", err)
		}

		returned_user, err := redis_client.GetUser(key)

		fmt.Println(returned_user, err)

		if expected_val != actual_val {
			t.Logf("Expected:\t %s \nGot:\t %s\n", expected_val, actual_val)
		}

	}
}
*/