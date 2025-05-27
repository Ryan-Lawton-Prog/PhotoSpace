package token

import (
	"sync"
)

type JWT string

type token struct {
	JWT  JWT
	Name string
}

var lock = &sync.Mutex{}

var instance *token

type newToken bool

func GetInstance() (*token, newToken) {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			instance = &token{}
			return instance, true
		}
	}

	return instance, false
}
