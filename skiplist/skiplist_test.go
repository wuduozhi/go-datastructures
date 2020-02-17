package skiplist

import (
	"testing"
)

func check(t *testing.T, s *SkipList, key, wanted int) {
	if got, _ := s.Get(key); got != wanted {
		t.Errorf("For key %v wanted value %v, got %v.", key, wanted, got)
	}
}

func TestSet(t *testing.T) {
	s := NewIntMap()
	if l := s.Len(); l != 0 {
		t.Errorf("Len is not 0, it is %v", l)
	}

	s.Set(0, 0)
	s.Set(1, 1)
	if l := s.Len(); l != 2 {
		t.Errorf("Len is not 2, it is %v", l)
	}
	check(t, s, 0, 0)
	if t.Failed() {
		t.Errorf("header.Next() after s.Set(0, 0) and s.Set(1, 1): %v.", s.header.next())
	}
	check(t, s, 1, 1)
}

func TestGet(t *testing.T) {
	s := NewIntMap()
	s.Set(0, 0)

	if value, present := s.Get(0); !(value == 0 && present) {
		t.Errorf("%v, %v instead of %v, %v", value, present, 0, true)
	}

	if value, present := s.Get(100); value != nil || present {
		t.Errorf("%v, %v instead of %v, %v", value, present, nil, false)
	}
}

func TestDelete(t *testing.T) {
	s := NewIntMap()
	for i := 0; i < 10; i++ {
		s.Set(i, i)
	}
	for i := 0; i < 10; i += 2 {
		s.Delete(i)
	}

	for i := 0; i < 10; i += 2 {
		if _, present := s.Get(i); present {
			t.Errorf("%d should not be present in s", i)
		}
	}

	if v, present := s.Delete(10000); v != nil || present {
		t.Errorf("Deleting a non-existent key should return nil, false, and not %v, %v.", v, present)
	}
}

func TestIter(t *testing.T) {
	s := NewIntMap()

	for i := 0; i < 20; i++ {
		s.Set(i, i)
	}

	seen := 0
	lastKey := 0

	i := s.Iterator()
	defer i.Close()

	for i.Next() {
		seen++
		lastKey = i.Key().(int)
		if i.Key() != i.Value() {
			t.Errorf("Wrong value for key %v: %v.", i.Key(), i.Value())
		}
	}

	if seen != s.Len() {
		t.Errorf("Not all the items in s where iterated through (seen %d, should have seen %d). Last one seen was %d.", seen, s.Len(), lastKey)
	}

	for i.Previous() {
		if i.Key() != i.Value() {
			t.Errorf("Wrong value for key %v: %v.", i.Key(), i.Value())
		}

		if i.Key().(int) >= lastKey {
			t.Errorf("Expected key to descend but ascended from %v to %v.", lastKey, i.Key())
		}

		lastKey = i.Key().(int)
	}

	if lastKey != 0 {
		t.Errorf("Expected to count back to zero, but stopped at key %v.", lastKey)
	}
}
