package main

import (
	"bytes"
	"crypto/md5"
	"io/ioutil"
	"net/http"
	"sort"
)

type Deleted int

const (
	NO   Deleted = 0
	SOFT Deleted = 1
	HARD Deleted = 2
)

// md5 for consistent hasing of services
type Record struct {
	rvolumes []string
	deleted  Deleted
	hash     string
}

type BodyReader struct {
	request *http.Request
}

// volume to sort by
type sortvol struct {
	score  []byte
	volume string
}
type byScore []sortvol

func (s byScore) Len() int      { return len(s) }
func (s byScore) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byScore) Less(i, j int) bool {
	return bytes.Compare(s[i].score, s[j].score) == 1
}

func key2voluem(key []byte, volumes []string, replicas int, subvolumes int) {
	var sortvols []sortvol
	for _, v := range volumes {
		hash := md5.New()
		hash.Write(key)
		hash.Write([]byte(v))
		sum := hash.Sum(nil)
		sortvols = append(sortvols, sortvol{sum, v})
	}
	// sort everything by score
	sort.Stable(byScore(sortvols))

}

func (br *BodyReader) readBody() []byte {
	body, err := ioutil.ReadAll(br.request.Body)
	if err != nil {
		panic("Error reading body")
	}
	return body
}
