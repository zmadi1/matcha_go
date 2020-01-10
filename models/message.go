package models

import (
	"fmt"

	"github.com/gmohlamo/matcha/database"
	"github.com/gmohlamo/matcha/mlogger"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Message struct {
	From    string `json:"from" bson:"from"`
	To      string `json:"to" bson:"to"`
	Message string `json:"message" bson:"message"`
}

type Chat struct {
	ID       bson.ObjectId          `json:"id" bson:"_id"`
	Users    map[string]interface{} `json:"users" bson:"users"`
	Messages []Message              `json:"messages" bson:"messages"`
}

var logger = mlogger.GetInstance()

func StoreMessage(message Message) {
	logger.Println("Storing message between: ", message.From, " and", message.To)
	client := database.GetInstance()
	defer client.Close()
	c := client.DB("matcha").C("chats")
	index := mgo.Index{
		Key:        []string{"users"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		logger.Println("Error ensuring index: ", err)
		return
	}
	chat := new(Chat)
	users := make(map[string]interface{})
	users[message.From] = FindUser("username", message.From)
	users[message.To] = FindUser("username", message.To)
	fmt.Println(users)
	err = c.Find(bson.M{
		"users": users,
	}).One(chat)
	if err != nil {
		logger.Println("Error retrieving chat struct: ", err)
	}
}
