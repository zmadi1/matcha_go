package database

import (
	"os"
	"sync"

	"github.com/blavkboy/matcha/mlogger"
	mgo "gopkg.in/mgo.v2"
)

//Decided to rewrite my whole approach to using the database

var once sync.Once
var session *mgo.Session

//InitDB will perform singleton pattern initialization on the
//database and return a session connection to the database.
//Recommended to only in the main function so you are able to
//close the session and use copies of the session with the
//defer session.Close() method call ensuring that the connection is
//properly terminated when the program closes.
func InitDB() (error, *mgo.Session) {
	mlogger := mlogger.GetInstance()

	var err error = nil
	once.Do(func() {
		session, err = mgo.Dial("mongodb://localhost:27017")
		if err != nil {
			mlogger.Println("Error: ", err)
			os.Exit(1)
		}
		session.SetMode(mgo.Monotonic, true)
	})
	if err != nil {
		return err, nil
	}
	mlogger.Println("mongodb session struct: ", session)
	return nil, session
}

func GetInstance() *mgo.Session {
	return session.Copy()
}
