package main

import (
	"fmt"
	"net/http"
)

func (a *App) handleGet(rw *http.ResponseWriter, rq *http.Request) {
	key := []byte(rq.URL.Path)
	w := *rw
	if !a.LockKey(key) {
		w.WriteHeader(404)
	} else {
		defer a.UnlockKey(key)
	}
	data, err := a.db.Get(key, nil)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.Write(data)
	w.WriteHeader(200)
	return
}
func (a *App) handlePut(rw *http.ResponseWriter, rq *http.Request) {
	key := []byte(rq.URL.Path)
	w := *rw
	values, err := rq.URL.Query()["value"]
	if !err || len(values[0]) < 1 {
		fmt.Println("URL Parameter 'value' is missing")
		w.WriteHeader(400)
		return
	}
	if !a.LockKey(key) {
		w.WriteHeader(404)
	} else {
		defer a.UnlockKey(key)
	}
	valuesBytes := []byte(values[0])
	a.db.Put(key, valuesBytes, nil)
	w.WriteHeader(200)
	return
}

func (a *App) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	switch rq.Method {
	case "GET":
		a.handleGet(&rw, rq)
		return
	case "PUT":
		a.handlePut(&rw, rq)
		return
	default:
		rw.WriteHeader(400)
	}
	return
}
