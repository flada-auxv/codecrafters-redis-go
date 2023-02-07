package main

import (
	"errors"
	"fmt"
	"time"
)

type Store struct {
	store map[string]*StoreValue
	now   func() time.Time
}

type StoreValue struct {
	expiredAt time.Time
	value     string
}

func NewStore(now func() time.Time) *Store {
	return &Store{
		store: make(map[string]*StoreValue),
		now:   now,
	}
}

func (s Store) Get(k string) (string, error) {
	v, ok := s.store[k]
	if !ok {
		return "", nil
	}

	if !v.expiredAt.IsZero() && s.now().After(v.expiredAt) {
		v.expiredAt = time.Time{}
		v.value = ""
		return "", nil
	}

	return v.value, nil
}

func (s Store) Set(k string, v string) error {
	s.store[k] = &StoreValue{value: v, expiredAt: time.Time{}}
	return nil
}

func (s Store) SetWithExpiration(k string, v string, milliSecsToExpire int) error {
	fmt.Println("SetWithExpiration")
	if milliSecsToExpire > 1_000_000_000 {
		return errors.New("ERROR: The milliSecsToExpire is too big.")
	}
	s.store[k] = &StoreValue{value: v, expiredAt: s.now().Add(time.Millisecond * time.Duration(milliSecsToExpire))}
	fmt.Println(s.store[k].value)
	fmt.Println(s.store[k].expiredAt)
	fmt.Println(time.Now())
	return nil
}
