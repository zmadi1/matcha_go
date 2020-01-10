package socket

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gmohlamo/matcha/mlogger"
	"github.com/gmohlamo/matcha/models"
	"github.com/gmohlamo/matcha/views/components"
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2/bson"
)

//Idea here is to use the connections map to keep track of all the users who are connected to the server.
//In the event that a user is online, a new connection should be established, the User's ID and a pointer
//to his/her connection is stored in the map. Should a user disconnect then we delete the connection.

//Connection will represent the way that connections are dealt with.
//When instanciated the models.User struct will be used to identify
//the user who's browser is sending a message
type Connection struct {
	User       *models.User
	Connection *websocket.Conn
}

type ProfileForm struct {
	Fname       string   `json:"fname"`
	Lname       string   `json:"lname"`
	Email       string   `json:"email"`
	Gender      string   `json:"gender"`
	Orientation string   `json:"orientation"`
	Interests   []string `json:"interests"`
}

type MessageReader struct {
	Type        string         `json:"type"`
	CommandType string         `json:"commandType"`
	Component   string         `json:"component"`
	Command     string         `json:"command"`
	Message     models.Message `json:"message"`
	Pform       ProfileForm    `json:"pform"`
}

type Response struct {
	Status    string `json:"status"`
	Column    string `json:"column"`
	Component string `json:"component"`
}

var UserConnections = make(map[bson.ObjectId]Connection)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleMessage(msg *MessageReader) {
	if strings.Compare("message", string(msg.Type)) == 0 {
		message := new(models.Message)
		fmt.Println(message)
		return
	}
}

func (msg *MessageReader) HandleCommand(connection *Connection, u *models.User) {
	mlogger := mlogger.GetInstance()
	if msg.CommandType == "profile" {
		if msg.Pform.Lname != "" {
			u.Lname = msg.Pform.Lname
		}
		if msg.Pform.Fname != "" {
			u.Fname = msg.Pform.Fname
		}
		if msg.Pform.Gender != "Select" {
			u.Sex = msg.Pform.Gender
		}
		if msg.Pform.Orientation != "Select" {
			u.Profile.Orientation = msg.Pform.Orientation
		}
		if len(msg.Pform.Interests) > 0 && msg.Pform.Interests[0] != "" {
			u.Profile.Interests = msg.Pform.Interests
		}
	} else if msg.CommandType == "propic" {
		u.Profile.Propic = msg.Command
	}
	err := u.UpdateUser()
	if err != nil {
		fmt.Println("Error updating user: ", err)
		mlogger.Println("Error updating user: ", err)
	} else {
		var b bytes.Buffer
		foo := bufio.NewWriter(&b)
		components.RenderProfileColumn(&b, u)
		err = foo.Flush()
		if err != nil {
			mlogger.Println("Error flushing contents: ", err)
		} else {
			mlogger.Println("Managed to get the component from the views package")
			resp := new(Response)
			resp.Status = "success"
			resp.Column = msg.Component
			resp.Component = b.String()
			fmt.Println(resp)
			buf, err := json.Marshal(resp)
			if err != nil {
				fmt.Println("Error marshalling struct into a json string: ", err)
				mlogger.Println("Error marshalling struct into a json string: ", err)
			}
			err = connection.Connection.WriteMessage(websocket.TextMessage, buf)
			if err != nil {
				mlogger.Println("Error writing to the connection: ", err)
				fmt.Println("Error: ", err)
			}
			if err != nil {
				mlogger.Println("Error closing io.Writer: ", err)
				fmt.Println("Error closing io.Writer: ", err)
			}
		}
		fmt.Println("Success")
		mlogger.Println("Successully updated user_no: ", u.ID)
	}
}

func (msg *MessageReader) EvalMsg(user *models.User, connection *Connection) {
	switch msg.Type {
	case "command":
		msg.HandleCommand(connection, user)
	}
}
