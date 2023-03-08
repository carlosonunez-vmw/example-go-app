package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type fakeCounter struct{ value int }
type badFakeCounter struct{ value int }

func (f *fakeCounter) Add(val int) error {
	f.value += val
	return nil
}

func (f *fakeCounter) Get() int {
	return f.value
}

func (f *badFakeCounter) Add(val int) error {
	return errors.New("sorry, but I failed!")
}

func (f *badFakeCounter) Get() int {
	return f.value
}

func newFakeCounter() *fakeCounter {
	return &fakeCounter{value: 0}
}

func newBadFakeCounter() *badFakeCounter {
	return &badFakeCounter{value: -100}
}

// I was too lazy to make a session-aware test; sorry :(
func TestAddEndpointPassWhenValidInt(t *testing.T) {
	wants := []int{10, 20, 30}
	for idx, valInt := range wants {
		initialValue := valInt
		if idx == 0 {
			initialValue = 0
		}
		want := initialValue + valInt
		val := strconv.Itoa(valInt)
		cookie := http.Cookie{Name: "incrementSessionsCookie", Value: "nil"}
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/add?by=%s", val), nil)
		req.AddCookie(&cookie)
		w := httptest.NewRecorder()
		incrementHandler(w, req)
		res := w.Result()
		defer res.Body.Close()
		resultJSON, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("Expected no errors, but got one: %s", err.Error())
			t.Fail()
		}
		var result counterResult
		if err := json.Unmarshal(resultJSON, &result); err != nil {
			t.Errorf("Expected no errors while reading response JSON, but got one: %s", err.Error())
			t.Fail()
		}
		if result.NewValue != valInt {
			t.Errorf("new value '%d' does not match expected new value '%d'",
				result.NewValue, want)
			t.Fail()
		}
	}
}
