package user

import (
	"sync"
)

type user struct {
	JWT  JWT
	Name string
}

var lock = &sync.Mutex{}

var instance *user

type newUser bool

type JWT string

func GetInstance() (*user, newUser) {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &user{}
			return instance, true
		}
	}

	return instance, false
}
