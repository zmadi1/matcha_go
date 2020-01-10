package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gmohlamo/matcha/mlogger"
	"github.com/gmohlamo/matcha/models"
	"golang.org/x/crypto/bcrypt"
)

var mySigningKey = []byte("The `jig is up")

//NewToken produces a new token
func NewToken(w http.ResponseWriter, r *http.Request) {
	logger := mlogger.GetInstance()
	if strings.Compare(r.Method, "POST") == 0 {
		w.Header().Set("Content-Type", "application/json")
		var user models.User
		var compare *models.User
		json.NewDecoder(r.Body).Decode(&user)
		// Create a new token object, specifying signing method and the claims
		// you would like it to contain.
		compare = models.FindUser("username", user.Username)
		pass1 := []byte(compare.Password)
		pass2 := []byte(user.Password)
		if compare == nil {
			w.Write([]byte("{\"success\": false}"))
			return
		} else if bcrypt.CompareHashAndPassword(pass1, pass2) != nil {
			w.Write([]byte("{\"success\": false}"))
			return
		}
		logger.Println("Password match for user: ", bcrypt.CompareHashAndPassword(pass1, pass2) == nil)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": compare.Username,
			"email":    compare.Email,
			"created":  time.Now().Unix(),
			"exp":      time.Now().AddDate(0, 0, 7).Unix(),
		})
		logger.Println("DEBUG INFO user_id: ", compare.ID)
		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString(mySigningKey)
		if err != nil {
			logger.Println("Error: ", err)
			w.Write([]byte("{\"success\": false}"))
			return
		}
		w.Write([]byte("{\"success\": true, \"token\": " + "\"" + tokenString + "\"" + "}"))
	}
}

func ConfirmUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		logger := mlogger.GetInstance()
		cookie, err := r.Cookie(authToken)
		if err != nil {
			failedAuth(w, r)
			return
		}
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				logger.Println("Error: Unexpected signing method: ", token.Header["alg"])
				return nil, fmt.Errorf("Error: Unexpected signing method: %v", token.Header["alg"])
			}

			return mySigningKey, nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			logger.Println(claims["user"])
		} else {
			logger.Println("Error: failed to authorize user")
			failedAuth(w, r)
			return
		}
		next(w, r)
	}
}

func GetCurrentUser(r *http.Request) *models.User {
	token := r.Header.Get("Authorization")
	return GetUserFromString(token)
}

func GetUserFromString(tokenString string) *models.User {
	mlogger := mlogger.GetInstance()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			mlogger.Println("Error: Unexpected signing method. ", token.Header["alg"])
			return nil, fmt.Errorf("Error: Unexpected signing method: %v", token.Header["alg"])
		}

		return mySigningKey, nil
	})
	if err != nil {
		mlogger.Println("Error parsing jwt string: ", err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := models.FindUser("email", claims["email"].(string))
		return user
	}
	return nil
}

func failedAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Authorization Failed")
}

const authToken = "AuthToken17286983217313"
