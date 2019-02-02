package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var collection *mongo.Collection

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "auth", Time: time.Now()})
		//log.Println(err)
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {

	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		//log.Println(err)
		return false
	}

	return true
}

// OpenAuthCollection open the collection of mongoDB called auth
func OpenAuthCollection() {
	collection = logger.DB.Collection("auth")
}

// CheckJwt if token is JWT and have time
func CheckJwt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_token := r.Header.Get("Authorization")
	if _token == "" {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "Authorization header missed with token value"}`))
		return
	}

	token, _ := jwt.Parse(_token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	// When using `Parse`, the result `Claims` would be a map.

	// In another way, you can decode token to your struct, which needs to satisfy `jwt.StandardClaims`
	user := JWT{}
	token, _ = jwt.ParseWithClaims(_token, &user, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if token.Valid {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok": "token right"}`))
		return
	}
	w.WriteHeader(403)
	_, _ = w.Write([]byte(`{"error": "token wrong"}`))
	return
}

// SignJWT function to create jwt token
func SignJWT(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = r.ParseForm()
	username := strings.Join(r.Form["username"], "")
	password := strings.Join(r.Form["password"], "")

	if username == "" || password == "" {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "input form is incomplete"}`))
		return
	}

	pointer := collection.FindOne(context.Background(), bson.D{{"username", username}})
	raw, err := pointer.DecodeBytes()
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "username or password wrong"}`))
		return
	}

	_username := raw.Lookup("username").StringValue()
	_password := raw.Lookup("password").StringValue()

	if username != _username || !ValidPassword(_password, password) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "username or password wrong"}`))
		return
	}

	token := createTokenString(User{Username: username, Password: password})
	w.WriteHeader(201)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))

	return

}

func createTokenString(user User) string {
	expireToken := time.Now().Add(time.Hour * 3).Unix() // token only valid 3 hours
	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &JWT{
		User: user,
		Age:  30,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	})

	// token -> string. Only server knows this secret
	_token, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		logger.SetSysLog(logger.SystemLog{Log: logger.Log{Event: "critical", Description: err.Error()},
			Pkg: "auth", Time: time.Now()})
		//log.Fatalln(err)
	}
	return _token
}

// RegisterUser function use in V1 router for sign up user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = r.ParseForm()
	username := strings.Join(r.Form["username"], "")
	password := strings.Join(r.Form["password"], "")

	if username == "" || password == "" {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error": "input form is incomplete"}`))
		return
	}

	pointer := collection.FindOne(context.Background(), bson.D{{"username", username}})
	_, err := pointer.DecodeBytes()

	if err != nil {

		_, _ = collection.InsertOne(context.Background(), User{Username: username, Password: SetPassword(password)})

		w.WriteHeader(201)
		_, _ = w.Write([]byte(`{"ok": "registered"}`))
		return
	}

	w.WriteHeader(400)
	_, _ = w.Write([]byte(`{"error": "username is taken"}`))
	return
}
