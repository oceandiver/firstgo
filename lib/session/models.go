package session

import (
	"crypto/rand"
	"github.com/dgrijalva/jwt-go"
	//"crypto/subtle"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/guregu/null"
	"io"
	"net/url"
	"time"
)

const (
	TokenLength int           = 32
	TtlDuration time.Duration = 20 * time.Minute
)

const (
	mySigningKey string = "who1Are1YOU4!"
)

type User struct {
	Id        uint64      `db:"id"`
	Email     string      `db:"email"`
	Token     string      `db:"token"`
	Ttl       time.Time   `db:"ttl"`
	OriginUrl null.String `db:"originurl"`
	OId       string      `db:"OId"`
}

type Event struct {
	Id        uint64    `db:"id"`
	Owner     string    `db:"owner"`
	Location  string    `db:"location"`
	Time      time.Time `db:"time"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	OId       string    `db:"OId"`
}

type Attendance struct {
	EventId   string    `db:"event_id"`
	UserId    string    `db:"user_id"`
	Attend    bool      `db:"attend"`
	UpdatedAt time.Time `db:"updated_at"`
}

// RefreshToken refreshes Ttl and Token for the User.
func (u *User) RefreshToken() error {

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Header["kid"] = "mySigningKey"
	token.Claims["oid"] = u.OId
	token.Claims["exp"] = time.Now().Add(TtlDuration).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)
	u.Token = tokenString

	return err

}

// IsValidToken returns a bool indicating that the User's current token hasn't
// expired and that the provided token is valid.
func (u *User) IsValidToken(myToken string) bool {

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey
	})

	if token.Valid {
		fmt.Println("You look nice today")
		if u.OId == token.Claims["oid"] {
			return true
		} else {
			return false
		}
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return false
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			return false
		} else {
			fmt.Println("Couldn't handle this token:", err)
			return false
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return false
	}

	return true
	//return subtle.ConstantTimeCompare([]byte(u.Token), []byte(token)) == 1
}

func (u *User) UpdateOriginUrl(originUrl *url.URL) error {
	var nsOrigin null.String
	if err := nsOrigin.Scan(originUrl.String()); err != nil {
		return err
	}

	u.OriginUrl = nsOrigin

	return nil
}
