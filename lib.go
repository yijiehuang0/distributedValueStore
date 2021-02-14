package main

import ()

type Deleted int

const (
	NO   Deleted = 0
	SOFT Deleted = 1
	HARD Deleted = 2
)

// use the md5 for some repartioning of the database
type Record struct {
	rvolumes []string
	deleted  Deleted
	hash     string
}
