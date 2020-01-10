package models

import (
	"encoding/json"
	"github.com/gmohlamo/matcha/database"
	"github.com/gmohlamo/matcha/mlogger"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
)

//The user file here will deal with all structures and
//methods that will work on the user and have to deal with
//the user's part in the database

//User struct will represent a user and will server
//as a storage mechanism for the server to keep track
//of user data.
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Username string        `json:"username" bson:"username"`
	Fname    string        `json:"fname" bson:"fname"`
	Lname    string        `json:"lname" bson:"lname"`
	Sex      string        `json:"sex" bson:"sex"`
	Email    string        `json:"email" bson:"email"`
	Password string        `json:"password" bson:"password"`
	Location GeoLocation   `json:"location" bson:"Location"`
	Profile  Profile       `json:"profile" bson:"profile"`
}

//constant for the cost
var Cost int = 6

//NewUser iss effectively how a new user is registered onto the
//system with their username, email and password. All other details
//are only necessary when the user model design is decided and users
//can provide more information about themselves.
func NewUser(user *User) *User {
	mlogger := mlogger.GetInstance()
	//Either initialize the database or get an instance of it
	client := database.GetInstance()
	defer client.Close()
	c := client.DB("matcha").C("users")
	index := mgo.Index{
		Key:        []string{"username", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		mlogger.Println("Error: ", err)
		panic(err)
	}
	mlogger.Println("Ensured Index")
	user.ID = bson.NewObjectId()
	err = c.Insert(&user)
	if err != nil {
		mlogger.Println("Error: ", err)
		return nil
	}
	mlogger.Println("Inserting User")
	user = FindUser("username", user.Username)
	user.Password = ""
	return user
}

//CheckUpdate will be used to check if the another user exists with matching
//credentials namely the username and email if those two are checked and values
//that were already filled in are not set to an empty string then the method
//will return true otherwise a false value is returned
func (user *User) CheckUpdate(updatedUser User) bool {
	if updatedUser.Username == "" || strings.Compare(updatedUser.Username, user.Username) != 0 {
		if strings.Compare(updatedUser.Username, user.Username) != 0 {
			if FindUser("username", updatedUser.Username) != nil {
				return false
			} else {
				return false
			}
		}
	} else if updatedUser.Fname == "" && user.Fname != "" {
		return false
	} else if updatedUser.Lname == "" && user.Lname != "" {
		return false
	} else if updatedUser.Sex == "" && user.Sex != "" {
		return false
	} else if updatedUser.Profile.Orientation == "" && user.Profile.Orientation != "" {
		return false
	} else if updatedUser.Profile.Confirmed == false && user.Profile.Confirmed == true {
		return false
	}
	return true
}

func (user *User) UpdateDiff(update User) {
	if strings.Compare(user.Username, update.Username) != 0 {
		user.Username = update.Username
	}
	if strings.Compare(user.Fname, update.Fname) != 0 {
		user.Fname = update.Fname
	}
	if strings.Compare(user.Lname, update.Lname) != 0 {
		user.Lname = update.Lname
	}
	if strings.Compare(user.Sex, update.Sex) != 0 {
		user.Sex = update.Sex
	}
	if strings.Compare(user.Profile.Orientation, update.Profile.Orientation) != 0 {
		user.Profile.Orientation = update.Profile.Orientation
	}
	user.Profile.Propic = update.Profile.Propic
	user.Location.Coordinates = update.Location.Coordinates
	user.Profile.Images = update.Profile.Images
	user.Profile.Interests = update.Profile.Interests
	user.UpdateUser()
}

func (user *User) UpdateUser() error {
	mlogger := mlogger.GetInstance()
	mlogger.Println("Attempting to update user_no: ", user.ID)
	client := database.GetInstance()
	defer client.Close()
	c := client.DB("matcha").C("users")
	return c.Update(bson.M{"_id": user.ID}, user)
}

//FindUser will return a User struct of the user being queried
//based on key value pair of the caller's choosing. Still needs
//to be tested extensively.
func FindUser(key string, value interface{}) *User {
	body := new(User)
	mlogger := mlogger.GetInstance()
	client := database.GetInstance()
	defer client.Close()
	c := client.DB("matcha").C("users")
	err := c.Find(bson.M{
		key: value,
	}).One(body)
	if err != nil {
		mlogger.Println("Error: ", err)
		return nil
	}
	return (body)
}

func FindMatch(user *User, w http.ResponseWriter) {
	var users []User
	mlogger := mlogger.GetInstance()
	client := database.GetInstance()
	defer client.Close()
	c := client.DB("matcha").C("users")
	err := c.Find(bson.M{"Location": bson.M{
		"$near": bson.M{
			"$geometry": bson.M{
				"type":        "Point",
				"coordinates": user.Location.Coordinates,
			},
			"$maxDistance": 500000,
		},
	},
	}).Select(bson.M{"username": 1, "Location": 1,
		"profile": 1,
		"sex":     1}).Skip(int(user.Profile.Index)).Limit(4).All(&users)
	if err != nil {
		mlogger.Println("Failed to obtain users")
		mlogger.Println("Error: ", err)
		return
	} else {
		users = filterUsers(user, users)
		mJson, _ := json.Marshal(users)
		w.Write(mJson)
		mlogger.Println("Sending:\n%+v\nObjects sent to client", users)
	}
}
