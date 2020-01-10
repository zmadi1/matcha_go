package models

import (
	"fmt"
	"github.com/gmohlamo/matcha/database"
	"github.com/gmohlamo/matcha/mlogger"
	"gopkg.in/mgo.v2/bson"
)

type Like struct {
	ID  bson.ObjectId `json:"id" bson:"_id"`
	Uid bson.ObjectId `json:"uid" bson:"uid"`
	Mid bson.ObjectId `json:"mid" bson:"mid"`
}

func likeUser(like *Like) {
	var user User
	mlogger := mlogger.GetInstance()
	client := database.GetInstance()
	defer client.Close()
	c := client.DB("matcha").C("users")
	err := c.Find(bson.M{"_id": like.Mid}).One(&user)
	if err != nil {
		mlogger.Println(err)
		return
	}
	user.Profile.Likes++
	err = c.Update(bson.M{"_id": user.ID}, user)
	if err != nil {
		mlogger.Println(err)
		return
	}
}

func findDupLike(like *Like, likes []Like) bool {
	itr := 0
	fmt.Println("Number of likes made by the sending user: ", len(likes))
	fmt.Printf("%+v\n", likes)
	for itr < len(likes) {
		if likes[itr].Mid == like.Mid {
			return true
		}
		itr++
	}
	return false
}

func AddLike(like *Like) bool {
	var likes []Like
	mlogger := mlogger.GetInstance()
	mlogger.Println("Processing like --> ", like)
	fmt.Println("Processing like --> ", like)
	client := database.GetInstance()
	defer client.Close()
	c := client.DB("matcha").C("likes")
	err := c.Find(bson.M{"uid": like.Uid}).All(&likes)
	if findDupLike(like, likes) == true {
		mlogger.Println("User has already liked this other user")
		return true
	} else { //because we failed to find the same like from before, we have to create a like
		err = c.Insert(bson.M{
			"uid": like.Uid,
			"mid": like.Mid,
		})
		if err != nil {
			mlogger.Println(err)
			fmt.Println(err)
			return false
		}
		fmt.Println("About to insert like into database")
		likeUser(like)
	}
	return true
}
