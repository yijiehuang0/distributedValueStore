package main

import (
	"flag"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strings"
	"sync"
	"time"
)

type App struct {
	appLock sync.Mutex
	db      *leveldb.DB
	locks   map[string]struct{}

	volumes    []string
	fallback   string
	replicas   int
	subvolumes int
	protect    bool
	md5sum     bool
	voltimeout time.Duration
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

func (a *App) GetRecord(key []byte) {
}
func (a *App) PutRecord(key []byte) {
}

func main() {
	fmt.Println("Just going to server some data")
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100

	port := flag.Int("port", 5000, "Port for server to listen to")
	pdb := flag.String("db", "/tmp/ejvaluedb", "Path to leveldb")
	fallback := flag.String("fallback", "", "Fallback server for missing keys")
	replicas := flag.Int("replicas", 3, "Amount of replicas to make of the dat")
	subvolumes := flag.Int("subvolumes", 10, "Amount of subvolumes, disks per machine")
	pvolumes := flag.String("volumes", "", "Volumes to use for storage, comma separated")
	protect := flag.Bool("protect", false, "Force UNLINK before DELETE")
	md5sum := flag.Bool("md5sum", true, "Calculate and store MD5 checksum of values")
	voltimeout := flag.Duration("volttimeout", 1*time.Second, "Volume servers must respond to GET/HEAD REQUESTS in this amount of time")
	flag.Parse()

	database, err := leveldb.OpenFile(*pdb, nil)
	volumes := strings.Split(*pvolumes, ",")

	defer database.Close()
	if err != nil {
		panic(err)
	}
	app := App{
		db:         database,
		locks:      make(map[string]struct{}),
		protect:    *protect,
		md5sum:     *md5sum,
		subvolumes: *subvolumes,
		replicas:   *replicas,
		fallback:   *fallback,
		volumes:    volumes,
		voltimeout: *voltimeout,
	}
	error := http.ListenAndServe(fmt.Sprint("localhost:", *port), &app)
	if error != nil {
		panic(error)
	}
}
