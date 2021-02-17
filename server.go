package main

import (
	"fmt"
	"net/http"
)

// routes:
// GET /{key} return value
// POST /{key} -d {value}
// GET /query&
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
func (a *App) handlePutPost(rw *http.ResponseWriter, rq *http.Request) {
	key := []byte(rq.URL.Path)
	w := *rw
	if rq.ContentLength == 0 {
		fmt.Println("Zero data in body")
		w.WriteHeader(400)
		return
	}
	if !a.LockKey(key) {
		w.WriteHeader(404)
	} else {
		defer a.UnlockKey(key)
	}
	bodyReader := BodyReader{request: rq}
	valuesBytes := bodyReader.readBody()
	a.db.Put(key, valuesBytes, nil)
	w.WriteHeader(200)
	return
}

func (a *App) handlePutPostDistributed(rw *http.ResponseWriter, rq *http.Request) {
	key := []byte(rq.URL.Path)
	w := *rw
	if rq.ContentLength == 0 {
		fmt.Println("Zero data in body")
		w.WriteHeader(400)
		return
	}
	if !a.LockKey(key) {
		w.WriteHeader(404)
	} else {
		defer a.UnlockKey(key)
	}
	bodyReader := BodyReader{request: rq}
	valuesBytes := bodyReader.readBody()
	a.db.Put(key, valuesBytes, nil)
	w.WriteHeader(200)
	return
}
func (a *App) handleQuery(rw *http.ResponseWriter, rq *http.Request) {
	w := *rw
	if len(rq.URL.RawQuery) > 0 {
		if rq.Method != "GET" {
			w.WriteHeader(403)
			return
		}
	}
}

func (a *App) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	a.handleQuery(&rw, rq)
	switch rq.Method {
	case "GET":
		a.handleGet(&rw, rq)
		return
	case "PUT", "POST":
		a.handlePutPost(&rw, rq)
		return
	default:
		rw.WriteHeader(400)
		return
	}
}
