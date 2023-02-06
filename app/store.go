package main

import (
	"errors"
	"sync"
	"time"
)

type Store struct {
	store map[string]StoreValue
	mu    *sync.Mutex
}

type StoreValue struct {
	expiredAt time.Time
	value     string
}

func (s Store) Get(k string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.store[k]
	if !ok {
		return "", errors.New("ERROR: Missing key.")
	}

	if time.Now().After(v.expiredAt) {
		s.store[k] = StoreValue{}
		return "", errors.New("ERROR: The key expired") // TODO: should it be returned with error??
	}

	return s.store[k].value, nil
}

func (s Store) Set(k string, v string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[k] = StoreValue{value: v, expiredAt: time.Time{}}
	return nil
}

func (s Store) SetWithExpiration(k string, v string, expiration int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if expiration > 1_000_000_000 {
		return errors.New("ERROR: The expiration is too big.")
	}
	s.store[k] = StoreValue{value: v, expiredAt: time.Now().Add(time.Second * time.Duration(expiration))}
	return nil
}
