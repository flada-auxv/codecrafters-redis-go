package main

import (
	"fmt"
	"testing"
	"time"
)

func fakeNow() time.Time {
	t, err := time.Parse("Jan 2, 2006 at 3:04pm", "Feb 7, 2023 at 1:00pm")
	if err != nil {
		panic(fmt.Sprintf("error: %#v", err))
	}
	return t
}

func TestStore_Set(t *testing.T) {
	t.Run("SET sets the v(alue) to the k(ey)", func(t *testing.T) {
		s := NewStore(time.Now)

		if gotErr := s.Set("testKey", "testValue"); gotErr != nil {
			t.Errorf("got error: %#v", gotErr)
		}
		if s.store["testKey"].value != "testValue" {
			t.Errorf("value is not set as expected")
		}
		if !s.store["testKey"].expiredAt.IsZero() {
			t.Errorf("expiredAt is not set as expected")
		}
	})
}

func TestStore_Get(t *testing.T) {
	t.Run("GET gets the v(alue) at the k(ey)", func(t *testing.T) {
		s := &Store{store: map[string]*StoreValue{"testKey": {value: "testValue"}}, now: fakeNow}
		v, err := s.Get("testKey")

		if err != nil {
			t.Errorf("got error: %#v", err)
		}
		if v != "testValue" {
			t.Errorf("value is not expected value. got: %#v", v)
		}
	})

	t.Run("when no value is set for the key", func(t *testing.T) {
		s := &Store{store: map[string]*StoreValue{"testKey": {value: "testValue"}}, now: fakeNow}
		v, err := s.Get("theKeyWithNoValueSet")

		if err != nil {
			t.Errorf("got error: %#v", err)
		}
		if v != "" {
			t.Errorf("value is not expected value. got: %#v", v)
		}
	})

	t.Run("when the value is expired", func(t *testing.T) {
		s := &Store{
			store: map[string]*StoreValue{"testKey": {value: "testValue", expiredAt: fakeNow().Add(time.Second * 10)}},
			now:   fakeNow,
		}
		v, err := s.Get("testKey")

		if err != nil {
			t.Errorf("got error: %#v", err)
		}
		if v != "testValue" {
			t.Errorf("value is not expected value. got: %#v", v)
		}

		s.now = func() time.Time {
			return fakeNow().Add(time.Second * 10)
		}
		v2, err2 := s.Get("testKey")

		if err2 != nil {
			t.Errorf("got error: %#v", err2)
		}
		if v2 != "testValue" {
			t.Errorf("value is not expected value. got: %#v", v2)
		}

		s.now = func() time.Time {
			return fakeNow().Add(time.Second * 11)
		}
		v3, err3 := s.Get("testKey")

		if err3 != nil {
			t.Errorf("got error: %#v", err)
		}
		if v3 != "" {
			t.Errorf("value is not expected value. got: %#v", v2)
		}
		if s.store["testKey"].value != "" {
			t.Errorf("internal value is not expected value. got: %#v", s.store["testKey"].value)
		}
	})
}

func TestStore_SetWithExpiration(t *testing.T) {
	t.Run("SET sets the v(alue) to the k(ey) with expiration", func(t *testing.T) {
		s := NewStore(time.Now)

		if gotErr := s.SetWithExpiration("testKey", "testValue", 10); gotErr != nil {
			t.Errorf("got error: %#v", gotErr)
		}
		if s.store["testKey"].value != "testValue" {
			t.Errorf("value is not set as expected")
		}
		if s.store["testKey"].expiredAt != fakeNow().Add(time.Second*10) {
			t.Errorf("expiredAt is not set as expected")
		}
	})
}
