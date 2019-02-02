package auth

import "github.com/dgrijalva/jwt-go"

// User model for auth staff
type User struct {
	Username string
	Password string
	Email    string
}

// JWT model for save in database
type JWT struct {
	User User
	Age  int
	jwt.StandardClaims
}

// SetPassword example usage of this two functions
// pass := SetPassword(password)
func SetPassword(pass string) string {
	return hashAndSalt([]byte(pass))
}

// ValidPassword bo := ValidPassword(pass, "as")
// fmt.Println(pass, bo)
func ValidPassword(hash string, pass string) bool {
	return comparePasswords(hash, []byte(pass))
}
