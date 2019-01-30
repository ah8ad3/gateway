package auth

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

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


// CheckJwt if token is JWT and have time
//func CheckJwt(w http.ResponseWriter, r *http.Request) {
//	_ = r.ParseForm()
//	tokenstring := strings.Join(r.Form["token"], "")
//	token, _ := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
//		return []byte(os.Getenv("SECRET_KEY")), nil
//	})
//	// When using `Parse`, the result `Claims` would be a map.
//
//	// In another way, you can decode token to your struct, which needs to satisfy `jwt.StandardClaims`
//	user := JWT{}
//	token, _ = jwt.ParseWithClaims(tokenstring, &user, func(token *jwt.Token) (interface{}, error) {
//		return []byte(os.Getenv("SECRET_KEY")), nil
//	})
//	if token.Valid {
//		_, _ = fmt.Fprintf(w, "this token is right")
//		return
//	}
//	_, _ = fmt.Fprintf(w, "this token is wrong")
//	return
//}

// SignJWT function to create jwt token
//func SignJWT(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		_ = r.ParseForm()
//		username := strings.Join(r.Form["username"], "")
//		password := strings.Join(r.Form["password"], "")
//
//		if username == "" || password == "" {
//			_, _ = fmt.Fprintf(w, "input form are incomplete")
//			return
//		}
//
//		user := User{}
//		DB.Where("username = ?", username).First(&user)
//
//		if user.ID == 0 || !ValidPassword(user.Password, password){
//			_, _ = fmt.Fprintf(w, "username or password wrong")
//			return
//		}
//		token := createTokenString(user)
//		_, _ = fmt.Fprintf(w, token)
//		return
//	}
//}

