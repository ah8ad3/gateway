package auth

import (
	"context"
	"fmt"
	"github.com/ah8ad3/gateway/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
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
		log.Println(err)
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

func OpenAuthCollection() {
	collection = logger.DB.Collection("auth")
}


// CheckJwt if token is JWT and have time
func CheckJwt(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	tokenstring := strings.Join(r.Form["token"], "")
	token, _ := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	// When using `Parse`, the result `Claims` would be a map.

	// In another way, you can decode token to your struct, which needs to satisfy `jwt.StandardClaims`
	user := JWT{}
	token, _ = jwt.ParseWithClaims(tokenstring, &user, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if token.Valid {
		_, _ = fmt.Fprintf(w, "this token is right")
		return
	}
	_, _ = fmt.Fprintf(w, "this token is wrong")
	return
}

// SignJWT function to create jwt token
func SignJWT(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		_ = r.ParseForm()
		username := strings.Join(r.Form["username"], "")
		password := strings.Join(r.Form["password"], "")

		if username == "" || password == "" {
			_, _ = fmt.Fprintf(w, "input form are incomplete")
			return
		}

		pointer := collection.FindOne(context.Background(), bson.D{{"username", username}})
		raw, err := pointer.DecodeBytes()
		if err != nil {
			log.Fatal(err)
		}

		if raw == nil {
			_, _ = fmt.Fprintf(w, "username or password wrong")
			return
		}

		_username := string(raw.Lookup("username").Value[:])
		_password := string(raw.Lookup("password").Value[:])

		if username != _username || !ValidPassword(_password, password) {
			_, _ = fmt.Fprintf(w, "username or password wrong")
			return
		}

		token := createTokenString(User{Username: username, Password: password})
		_, _ = fmt.Fprintf(w, token)
		return
	}
}

func createTokenString(user User) string {
	expireToken := time.Now().Add(time.Hour * 24).Unix()
	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &JWT{
		User: user,
		Age:  30,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireToken,
		},
	})

	// token -> string. Only server knows this secret (foobar).
	tokenstring, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Fatalln(err)
	}
	return tokenstring
}



func ShowLogs() {
	// bson.D{{}} return all records in database you can create a map here
	//logg := logger.UserLog{}
	data1, _ := logger.Collection.Find(context.Background(), bson.D{{"log.event", "critical"}})

	for data1.Next(context.Background()) {
		raw, err := data1.DecodeBytes()
		if err != nil { log.Fatal(err) }
		logg := logger.Log{}
		_ = raw.Lookup("log").Unmarshal(&logg)
		fmt.Println(logg)
	}
}

