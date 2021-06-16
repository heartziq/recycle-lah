package utility

import (
	"io"
	"net/http"
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
func xVerifyToken(tokenString string) (bool, error) {
	// Verify
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})

	if err != nil {

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token has expired
				log.Println("token expired")
				return false, errors.New("token expired")
			}
		}

		log.Println("invalid token")
		return false, errors.New("invalid token")

	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Info.Println("Welcome, ", claims["aud"]) // claims["aud"] will hold the userid
		Info.Println("Token is valid.")
		return true, nil
	}

	return false, errors.New("Unknown error")

}

// VerifyToken checks if provided token is valid -> return true if valid
// error can be: "token expired", "invalid token"
// anything else will throw error unknown
func VerifyToken(tokenString string) (string, error) {
	Trace.Println("in verifyToken")
	// Verify
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})

	if err != nil {
		Trace.Println("there is error ", err)
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

	Trace.Printf("dynamic type: %T\n", token.Claims)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userid := claims["aud"]
		Trace.Println("userid=", userid)
		if v, ok := userid.(string); ok {
			log.Println("Token is valid.")
			return v, nil
		}

		// return "", errors.New("invalid token")
	}
	Trace.Println("after claims (error)")
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

// returns []byte hashed password on a given string
func HashPassword(password string) []byte {
	if hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost); err != nil {
		Error.Println(err)
		return nil
	} else {
		return hash
	}
}

// to retrieve user id from token in request header
// not working
func GetUserId(r *http.Request) (string, error) {
	Trace.Println("============ get user id  =====================")
	mechanism := strings.Split(r.Header.Get("Authorization"), " ")
	if len(mechanism) > 1 && mechanism[0] == "Bearer" {
		if token := mechanism[1]; token != "" {
			userId, err := VerifyToken(token)

			// validate token
			if err != nil {
				Trace.Println(err)
				return "", err
			}
			Trace.Println(userId)
			return userId, nil
		}

	}
	return "", errors.New("error getting userid from token")
}
