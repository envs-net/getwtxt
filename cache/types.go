package cache

import (
	"sync"
)

// Indexer allows for other uses of the Index functions
type Indexer interface {
	AddUser(string, string)
	DelUser(string)
}

// UserIndex provides an index of users by URL
type UserIndex map[string]*Data

// Data from user's twtxt.txt
type Data struct {
	nick   string
	url    string
	status []string
}

// Mutex to control access to the User Index
var imutex = sync.RWMutex{}