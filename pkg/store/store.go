package store

import (
	"errors"
	"time"
)

type Store interface {
	Get(k string) (string, error)
	Set(k string, v string) error
	SetWithExpiration(k string, v string, milliSecsToExpire int) error
}

type MemoryStore struct {
	store map[string]*MemoryStoreValue
	now   func() time.Time
}

type MemoryStoreValue struct {
	expiredAt time.Time
	value     string
}

func NewMemoryStore(now func() time.Time) *MemoryStore {
	return &MemoryStore{
		store: make(map[string]*MemoryStoreValue),
		now:   now,
	}
}

func (s MemoryStore) Get(k string) (string, error) {
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

func (s MemoryStore) Set(k string, v string) error {
	s.store[k] = &MemoryStoreValue{value: v, expiredAt: time.Time{}}
	return nil
}

func (s MemoryStore) SetWithExpiration(k string, v string, milliSecsToExpire int) error {
	if milliSecsToExpire > 1_000_000_000 {
		return errors.New("ERROR: The milliSecsToExpire is too big")
	}
	s.store[k] = &MemoryStoreValue{value: v, expiredAt: s.now().Add(time.Millisecond * time.Duration(milliSecsToExpire))}
	return nil
}
