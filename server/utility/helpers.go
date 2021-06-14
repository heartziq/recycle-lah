package utility

import (
	"io"
	"os"
	"strings"
	"time"

	"errors"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	SECRET_FILE = "secret.pem"
	KEY         string
)

func init() {
	KEY = initSecretKey() // Secret key to sign jwt

}

func initSecretKey() string {

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	file, err := os.Open(path + "\\utility\\" + SECRET_FILE)
	if err != nil {
		panic(err)
	}

	content, _ := io.ReadAll(file)

	secret := strings.TrimSpace(string(content))
	secret = strings.TrimPrefix(string(secret), "-----BEGIN RSA PRIVATE KEY-----")
	secret = strings.TrimSuffix(string(secret), "-----END RSA PRIVATE KEY-----")
	return secret
}

// GenToken generates jwt token
func GenToken(secret, userid string) (string, error) {

	mySigningKey := []byte(secret)
	expiryDate := time.Now().Add(time.Hour * 24 * 7).Unix()

	// get userid

	// Create the Claims
	claims := &jwt.StandardClaims{
		Audience:  userid,
		ExpiresAt: expiryDate,
		Issuer:    "test",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}

	return ss, nil
}

// VerifyToken checks if provided token is valid -> return true if valid
// error can be: "token expired", "invalid token"
// anything else will throw error unknown
func VerifyToken(tokenString string) (string, error) {
	// Verify
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})

	if err != nil {

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token has expired
				log.Println("token expired")
				return "", errors.New("token expired")
			}
		}

		log.Println("invalid token")
		return "", errors.New("invalid token")

	}

	// fmt.Printf("dynamic type: %T\n", token.Claims)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userid := claims["aud"]
		if v, ok := userid.(string); ok {
			log.Println("Token is valid.")
			return v, nil
		}

		// return "", errors.New("invalid token")
	}

	return "", errors.New("Unknown error")

}

// VerifyPassword checks if password is correct
func VerifyPassword(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}
