package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"sync"
)

type App struct {
	appLock sync.Mutex
	db      *leveldb.DB
	locks   map[string]struct{}
}

func (a *App) LockKey(input []byte) bool {
	a.appLock.Lock()
	defer a.appLock.Unlock()
	s := string(input)
	if _, contains := a.locks[s]; contains { // someone else is holding this lock
		return false
	}
	a.locks[s] = struct{}{}
	return true
}

func (a *App) UnlockKey(input []byte) {
	a.appLock.Lock()
	defer a.appLock.Unlock()
	delete(a.locks, string(input))
}

func main() {
	// going to define the max connections per hose
	fmt.Println("Just going to server some data")
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	database, err := leveldb.OpenFile("/tmp/ejvaluedb", nil)
	defer database.Close()
	if err != nil {
		panic(err)
	}
	app := App{
		db:    database,
		locks: make(map[string]struct{}),
	}
	error := http.ListenAndServe("localhost:5000", &app)
	if error != nil {
		panic(error)
	}
}
