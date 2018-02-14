package user

import (
	"fmt"
	"manageiq-exchange/models/utils"
	"reflect"
	"testing"
)

func TestUserInit(t *testing.T) {

	want := User{
		Login:    "aljesusg",
		Name:     "Alberto",
		GithubID: 1,
	}
	var data = map[string]interface{}{
		"login":     "aljesusg",
		"github_id": 1,
		"name":      "Alberto",
	}
	var user User
	user.Init(data)
	if !reflect.DeepEqual(user, want) {
		t.Errorf("UserInit returned %+v, want %+v", user, want)
	}
}

func TestUserPrint(t *testing.T) {
	want := fmt.Sprintf("%s: aljesusg (1)\n\n", utils.PrintColor("User", "Red"))
	want += fmt.Sprintf("    Name : Alberto\n")

	user := User{
		Login:    "aljesusg",
		Name:     "Alberto",
		GithubID: 1,
	}
	if !reflect.DeepEqual(user.Print(), want) {
		t.Errorf("User Print returned -%+v-, want -%+v-", user.Print(), want)
	}
}

func TestUserCollectionInit(t *testing.T) {
	user := User{
		Login:    "aljesusg",
		Name:     "Alberto",
		GithubID: 1,
	}
	users := UserCollection{}
	users.Users = append(users.Users, user)
	users.Total = len(users.Users)

	var data = map[string]interface{}{
		"login":     "aljesusg",
		"github_id": 1,
		"name":      "Alberto",
	}

	var dataUsers = []interface{}{}
	dataUsers = append(dataUsers, data)

	var userCollection UserCollection
	userCollection.Init(dataUsers)

	if !reflect.DeepEqual(userCollection, users) {
		t.Errorf("UserCollectionInit returned %+v, want %+v", userCollection, users)
	}
}

func TestUserCollectionPrint(t *testing.T) {
	want := fmt.Sprintf("\n\n%s (%d/%d)\n\n------------------------------------------\n\n", utils.PrintColor("Users", "Red"), 1, 1)
	want += fmt.Sprintf("%s: aljesusg (1)\n\n", utils.PrintColor("User", "Red"))
	want += fmt.Sprintf("    Name : Alberto\n\n\n")

	user := User{
		Login:    "aljesusg",
		Name:     "Alberto",
		GithubID: 1,
	}
	userCollection := UserCollection{}
	userCollection.Users = append(userCollection.Users, user)
	userCollection.Total = len(userCollection.Users)
	gotUserCollection := userCollection.Print(1)

	if !reflect.DeepEqual(gotUserCollection, want) {
		t.Errorf("UserCollection Print returned -%+v-, want -%+v-", gotUserCollection, want)
	}
}
