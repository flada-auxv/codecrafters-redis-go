package main

import (
	"errors"
	"time"
)

type Store struct {
	store map[string]StoreValue
}

type StoreValue struct {
	expiredAt time.Time
	value     string
}

func NewStore() *Store {
	return &Store{
		store: make(map[string]StoreValue),
	}
}

func (s Store) Get(k string) (string, error) {
	v, ok := s.store[k]
	if !ok {
		return "", errors.New("ERROR: Missing key.")
	}

	if !v.expiredAt.IsZero() && time.Now().After(v.expiredAt) {
		v = StoreValue{}
		return "", errors.New("ERROR: The key expired") // TODO: should it be returned with error??
	}

	return v.value, nil
}

func (s Store) Set(k string, v string) error {
	s.store[k] = StoreValue{value: v, expiredAt: time.Time{}}
	return nil
}

func (s Store) SetWithExpiration(k string, v string, expiration int) error {
	if expiration > 1_000_000_000 {
		return errors.New("ERROR: The expiration is too big.")
	}
	s.store[k] = StoreValue{value: v, expiredAt: time.Now().Add(time.Second * time.Duration(expiration))}
	return nil
}
