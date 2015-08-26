package session

import (
	"crypto/rand"
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
	Owner     uint64    `db:"owner"`
	Location  string    `db:"location"`
	Time      time.Time `db:"time"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	OId       string    `db:"OId"`
}

type Attendance struct {
	EventId   uint64    `db:"event_id"`
	UserId    uint64    `db:"user_id"`
	Attend    bool      `db:"attend"`
	UpdatedAt time.Time `db:"updated_at"`
}

// RefreshToken refreshes Ttl and Token for the User.
func (u *User) RefreshToken() error {
	t := time.Now().UTC().Add(TtlDuration)
	gobEncoded, err := t.GobEncode()

	if err != nil {
		fmt.Println(err)
	}

	token := make([]byte, TokenLength)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		return err
	}

	fmt.Println(gobEncoded, token)
	c := [][]byte{gobEncoded, token}

	u.Token = base64.URLEncoding.EncodeToString(bytes.Join(c, []byte(", ")))
	u.Ttl = time.Now().UTC().Add(TtlDuration)
	return nil
}

// IsValidToken returns a bool indicating that the User's current token hasn't
// expired and that the provided token is valid.
func (u *User) IsValidToken(token string) bool {

	b, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(b)

	fmt.Println("---------------")

	c := bytes.Split(b, []byte(", "))
	fmt.Println(c)

	var t time.Time
	t.GobDecode(c[0])

	if t.Before(time.Now().UTC()) {
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
