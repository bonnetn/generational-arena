package arena

import (
	"testing"
)

func TestIterator_Next(t *testing.T) {
	var a Arena[int]
	a.Insert(1)
	a.Insert(2)
	a.Insert(3)

	it := a.Iterator()
	gotValue, gotId, done := it.Next()
	if done {
		t.Errorf("Next() returned done too early")
	}
	if gotValue != 1 {
		t.Errorf("Next() returned invalid value")
	}
	expectedID := ID{0, 1}
	if gotId != expectedID {
		t.Errorf("Next() returned invalid id")
	}

	gotValue, gotId, done = it.Next()
	if done {
		t.Errorf("Next() returned done too early")
	}
	if gotValue != 2 {
		t.Errorf("Next() returned invalid value")
	}
	expectedID = ID{1, 1}
	if gotId != expectedID {
		t.Errorf("Next() returned invalid id")
	}

	gotValue, gotId, done = it.Next()
	if done {
		t.Errorf("Next() returned done too early")
	}
	if gotValue != 3 {
		t.Errorf("Next() returned invalid value")
	}
	expectedID = ID{2, 1}
	if gotId != expectedID {
		t.Errorf("Next() returned invalid id")
	}

	_, _, done = it.Next()
	if !done {
		t.Errorf("Next() did not return done")
	}
}
