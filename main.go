package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
)

type Counter interface {
	Add(v int) error
	Get() int
}

type counter struct {
	value int
}

type counterResult struct {
	InitialValue int `json:"initial"`
	NewValue     int `json:"new"`
}

type counterError struct {
	Message string `json:"message"`
}

func (c *counter) Add(v int) error {
	c.value += v
	return nil
}

func (c *counter) Get() int {
	return c.value
}

const (
	cookieName = "incrementSessionsCookie"
)

var (
	key        = []byte("123456789012345")
	store      = sessions.NewCookieStore(key)
	counterMap = map[int64]*counter{}
)

func incrementError(w http.ResponseWriter, msg string) {
	e, err := json.Marshal(&counterError{Message: msg})
	if err != nil {
		http.Error(w, "failed to create error object", http.StatusInternalServerError)
	}
	http.Error(w, string(e), http.StatusUnprocessableEntity)
}

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToLower(r.Method) != "put" {
		incrementError(w, "please resend as a PUT")
		return
	}
	vals, ok := r.URL.Query()["by"]
	if !ok {
		incrementError(w, "please provide 'by'")
		return
	}
	val, err := strconv.Atoi(vals[0])
	if err != nil {
		incrementError(w, fmt.Sprintf("not a number: %s", vals[0]))
		return
	}
	var c *counter
	session, _ := store.Get(r, cookieName)
	log.Printf("session: %+v", session.Values)
	sessionID, ok := session.Values["sessionid"]
	if !ok {
		id := time.Now().Unix()
		log.Printf("saving new sessionid %d", id)
		session.Values["sessionid"] = fmt.Sprintf("%d", id)
		c = &counter{value: 0}
		counterMap[id] = c
		sessionID = id
	} else {
		idStr, ok := sessionID.(string)
		if !ok {
			http.Error(w, fmt.Sprintf("unable to parse session ID: %+v", sessionID), http.StatusInternalServerError)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid session ID: %s", idStr), http.StatusInternalServerError)
			return
		}
		if _, ok := counterMap[id]; !ok {
			log.Printf("initializing session ID %s", idStr)
			counterMap[id] = &counter{value: 0}
		}
		c = counterMap[id]
		sessionID = id
	}
	res := counterResult{InitialValue: c.value}
	log.Printf("[%d] before add: %d", sessionID, &c)
	if err := c.Add(val); err != nil {
		incrementError(w, fmt.Sprintf("increment failed: %s", err.Error()))
		return
	}
	log.Printf("[%d] before add: %d", sessionID, c)
	res.NewValue = c.value
	out, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, "failed to render incremented object", http.StatusInternalServerError)
	}
	session.Save(r, w)
	fmt.Fprintf(w, string(out))
}

func main() {
	http.HandleFunc("/add", incrementHandler)
	log.Print("Example Go App is running on port 5000")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
