package session

import (
	//"os"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/url"
	"os/exec"
	"strconv"
	"time"
)

type newUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type newEvent struct {
	Owner    string    `json:"owner"`
	Location string    `json:"location"`
	Time     time.Time `json:"time"`
}

type newAttendance struct {
	EventId string `json:"eid"`
	UserId  string `json:"uid"`
}

const longForm = "Jan 2, 2006 at 3:04pm (MST)"

/*
 * Send an email to the user with login token.
 *
 */
func SendLoginEmail(scheme, host, email, token, uid string) {

	// Build login url
	params := url.Values{}
	params.Add("Token", token)
	params.Add("UserId", uid)

	loginUrl := url.URL{}

	loginUrl.Scheme = scheme
	loginUrl.Host = host

	loginUrl.Path = "/verify"

	fmt.Println(token)
	fmt.Println(uid)
	fmt.Println(loginUrl)
	// Send login email
	var mailContent bytes.Buffer
	ctx := struct {
		LoginUrl string
	}{
		fmt.Sprintf("%s?%s", loginUrl.String(), params.Encode()),
	}

	go func() {
		if err := loginEmailTemplate.Execute(&mailContent, ctx); err == nil {

			//fmt.Fprintf(w, email)
			if err := SendMail([]string{email}, "Passwordless Login Verification", mailContent.String()); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Error sending verification email")
			}
		}
	}()
}

/*
 * curl  -D - -H "Content-Type: application/json" -d  \
 *   '{"name":"Joe", "email":"zhiruo@yahoo.com"}'     \
 *   http://130.211.143.134:8080/v1/user/signup
 */
func AddUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("in AddUser")
	var user newUser
	/*
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	*/
	user = newUser{}
	user.Name = r.FormValue("name")
	user.Email = r.FormValue("email")

	fmt.Println(user)

	/* write to db

	*/
	nUser := &User{Email: user.Email}
	email := user.Email

	if email != "" {

		if err := dbmap.SelectOne(nUser, "SELECT * FROM users WHERE email=$1", email); err != nil {
			fmt.Println(email)
			log.WithFields(log.Fields{
				"error": err,
				"user":  user,
			}).Warn("Error finding User.")
			nUser.Email = email
			/* generat a uuid */
			out, err := exec.Command("uuidgen").Output()
			if err != nil {
				log.Fatal(err)
				return
			}
			nUser.OId = string(out)
			if err := dbmap.Insert(nUser); err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"user":  user,
				}).Warn("Error creating User.")
				return
			}
		} else {
			nUser.Email = email
			dbmap.Update(nUser)
		}

		fmt.Println("------------", nUser)
		//Send Email

		var scheme string
		var host string
		if r.URL.IsAbs() {
			scheme = r.URL.Scheme
			host = r.URL.Host
		} else {
			scheme = "http"
			host = r.Host
		}

		nUser.RefreshToken()
		fmt.Println(nUser)
		token := nUser.Token
		uid := strconv.FormatUint(nUser.Id, 10)

		SendLoginEmail(scheme, host, email, token, uid)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}

	return
}

// GetContextUser returns the User for the given request context or nil.
func GetAuthenticatedUser(r *http.Request) *User {

	user := &User{}
	user.OId = r.FormValue("UserId")
	user.Token = r.FormValue("Token")

	if len(user.OId) == 0 || len(user.Token) == 0 {
		return nil
	}

	if user.IsValidToken(user.Token) != true {
		return nil
	}

	return user
}

/*
 * Alternative to login by providing the email and password.
 * Response contains a login token and uid.
 */
func GetToken(w http.ResponseWriter, r *http.Request) {
	// Redirect logged in users
	fmt.Println("in GetUser")

	fmt.Println(r)
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Println("password==", password, ",email=", email)

	user := &User{}
	if err := dbmap.SelectOne(user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"user":  user,
		}).Warn("Error retrieving User.")
		return
	}

	fmt.Println(user)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	user.RefreshToken()
	fmt.Println(user)
	fmt.Println(user.Token)

	if err := json.NewEncoder(w).Encode(user.Token); err != nil {
		panic(err)
	}

	return
}

/*
 * curl  -D -  http://130.211.143.134:8080/v1/user?id=1
 * signin user and return the user profile data
 */
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Redirect logged in users
	fmt.Println("in GetUser")

	user := GetAuthenticatedUser(r)
	if user == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := dbmap.SelectOne(user, "SELECT * FROM users WHERE \"OId\"=$1", user.OId); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"user":  user,
		}).Warn("Error retrieving Event.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}

	return
}

/*
 * curl  -D -  http://130.211.143.134:8080/v1/user?id=1
 */
func GetUsers(w http.ResponseWriter, r *http.Request) {

	fmt.Println("in GetUsers")

	user := GetAuthenticatedUser(r)

	if user == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := dbmap.SelectOne(user, "SELECT * FROM users WHERE \"OId\"=$1", user.OId); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"user":  user,
		}).Warn("Error retrieving Event.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	/* verify if the user has enough permission */
	if user.Id != 0 {

		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}
	} else {
		var users []User
		if _, err := dbmap.Select(users, "SELECT * FROM users"); err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"users": users,
			}).Warn("Error retrieving Event.")
			return
		}
		if err := json.NewEncoder(w).Encode(users); err != nil {
			panic(err)
		}
	}

	return

}

/*
 * curl  -D - -H "Content-Type: application/json" -d '{"owner":1, "location":"beijing", "time":"0001-01-01T00:00:00Z"}' http://130.211.143.134:8080/v1/addevent
 */
func AddEvent(w http.ResponseWriter, r *http.Request) {

	fmt.Println("in AddEvent")

	user := GetAuthenticatedUser(r)
	if user == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("user=", user)

	var ev newEvent
	var err error
	ev.Owner = user.OId
	ev.Location = r.FormValue("Location")

	ev.Time, err = time.Parse(longForm, r.FormValue("Time"))
	if err != nil {
		fmt.Println("Wrong time format", r.FormValue("Time"))
		return
	}
	fmt.Println(ev)

	/* write to db
	   Id        int64       `db:"id"`
	   Owner     int64       `db:"owner"`
	   Location  string      `db:"location"`
	   Time      time.Time   `db:"time"`
	   CreatedAt time.Time   `db:"created_at"`
	   UpdatedAt time.Time   `db:"updated_at"`
	*/
	/* generat a uuid */
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
		return
	}

	nEv := &Event{Owner: ev.Owner, OId: string(out), Location: ev.Location, Time: ev.Time, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if err := dbmap.Insert(nEv); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ev":    ev,
		}).Warn("Error creating Event.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ev); err != nil {
		panic(err)
	}

	return
}

/*
 * curl  -D - -H "Content-Type: application/json" -d '{"owner":1, "location":"beijing", "time":"0001-01-01T00:00:00Z"}' http://130.211.143.134:8080/v1/addevent
 */
func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in DeleteEvent")

	user := GetAuthenticatedUser(r)
	if user == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("user=", user)

	var ev newEvent
	var err error
	ev.Owner = user.OId
	ev.Location = r.FormValue("location")

	ev.Time, err = time.Parse(longForm, r.FormValue("Time"))
	if err != nil {
		fmt.Println("Wrong time format", r.FormValue("Time"))
		return
	}

	fmt.Println(ev)

	/* write to db
	   Id        int64       `db:"id"`
	   Owner     int64       `db:"owner"`
	   Location  string      `db:"location"`
	   Time      time.Time   `db:"time"`
	   CreatedAt time.Time   `db:"created_at"`
	   UpdatedAt time.Time   `db:"updated_at"`
	*/
	nEv := &Event{}

	if _, err := dbmap.Delete(nEv); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ev":    ev,
		}).Warn("Error creating Event.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ev); err != nil {
		panic(err)
	}
}

/*
 * curl  -D - -H "Content-Type: application/json" -d '{"owner":1, "location":"beijing", "time":"0001-01-01T00:00:00Z"}' http://130.211.143.134:8080/v1/addevent
 */
func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in Update Event")

	user := GetAuthenticatedUser(r)
	if user == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("user=", user)

	var ev newEvent
	var err error
	ev.Owner = user.OId
	ev.Location = r.FormValue("location")
	ev.Time, err = time.Parse(longForm, r.FormValue("Time"))
	if err != nil {
		fmt.Println("Wrong time format", r.FormValue("Time"))
		return
	}

	fmt.Println(ev)

	/* write to db
	   Id        int64       `db:"id"`
	   Owner     int64       `db:"owner"`
	   Location  string      `db:"location"`
	   Time      time.Time   `db:"time"`
	   CreatedAt time.Time   `db:"created_at"`
	   UpdatedAt time.Time   `db:"updated_at"`
	*/
	nEv := &Event{Owner: ev.Owner, Location: ev.Location, Time: ev.Time, CreatedAt: time.Now(), UpdatedAt: time.Now()}

	if _, err := dbmap.Update(nEv); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ev":    ev,
		}).Warn("Error creating Event.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ev); err != nil {
		panic(err)
	}

}

/*
 * curl  -D -  http://130.211.143.134:8080/v1/event?id=1
 */
func GetEvent(w http.ResponseWriter, r *http.Request) {
	// Redirect logged in users
	fmt.Println("in GetEvent")

	user := GetAuthenticatedUser(r)
	if user == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("user=", user)

	/* which event */
	eid := r.FormValue("EventId")

	fmt.Println("eid=", eid)

	/* write to db
	   Id        int64       `db:"id"`
	   Owner     int64       `db:"owner"`
	   Location  string      `db:"location"`
	   Time      time.Time   `db:"time"`
	   CreatedAt time.Time   `db:"created_at"`
	   UpdatedAt time.Time   `db:"updated_at"`
	*/
	nEv := &Event{}
	if err := dbmap.SelectOne(nEv, "SELECT * FROM events WHERE \"OId\"=$1", eid); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"ev":    nEv,
		}).Warn("Error retrieving Event.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nEv); err != nil {
		panic(err)
	}

	return
}

/*
 * curl  -D - -H "Content-Type: application/json" -d '{"EventId":1, "UserId":1 }' http://130.211.143.134:8080/v1/addattend
 */
func AddAttendance(w http.ResponseWriter, r *http.Request) {

	// TODO: check if the request is coming to the true user and relavent checks

	fmt.Println("in AddInvitee")

	user := GetAuthenticatedUser(r)
	if user == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("user=", user)

	eid := r.FormValue("EventId")
	uid := user.OId

	/* write to db
	   EventId       uint64       `db:"event_id"`
	   UserId        uint64       `db:"user_id"`
	   Attend        bool        `db:"attend"`
	   UpdatedAt     time.Time   `db:"updated_at"`
	*/
	nEv := &Attendance{EventId: eid, UserId: uid, Attend: false, UpdatedAt: time.Now()}

	if err := dbmap.Insert(nEv); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"nEv":   nEv,
		}).Warn("Error creating Event.")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	return
}

/*
 * curl  -D - -H "Content-Type: application/json" -d '{"EventId":1, "UserId":1 }' http://130.211.143.134:8080/v1/addattend
 */
func DeleteAttendance(w http.ResponseWriter, r *http.Request) {
}

/*
 * curl  -D -  http://130.211.143.134:8080/v1/user?id=1
 */
func GetEvents(w http.ResponseWriter, r *http.Request) {
}
