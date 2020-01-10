package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gmohlamo/matcha/database"
	"github.com/gmohlamo/matcha/mlogger"
	"github.com/gmohlamo/matcha/routing"
	"github.com/gmohlamo/matcha/services/auth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//In main we will handle all requests to the server
func main() {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	mlogger := mlogger.GetInstance()
	mlogger.Println(time.Now())
	err, conn := database.InitDB()
	defer conn.Close()
	if err != nil {
		mlogger.Println("Error: ", err)
		os.Exit(1)
	}
	mlogger.Println("Starting 'Matcha' dating service API")
	r := mux.NewRouter()

	//Register users to the users collectioon in the matcha database
	//rebasing this code to make an api for the front end
	//r.HandleFunc("/ws/{token}", routing.SocketConn)
	r.HandleFunc("/users/update", routing.HandleUpdate).Methods("POST")
	r.HandleFunc("/users/login", auth.NewToken).Methods("POST")
	r.HandleFunc("/user", routing.HandleUser).Methods("POST", "GET")
	r.HandleFunc("/users/check", routing.HandleCheck).Methods("GET")
	r.HandleFunc("/users/matches", routing.HandleMatches).Methods("GET")
	r.HandleFunc("/users/matches/like", routing.HandleLikes).Methods("POST")
	http.ListenAndServe(":4040", handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins)(r))
}
