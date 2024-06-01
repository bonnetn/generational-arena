package arena

import "testing"

func TestNewArenaWithCapacity(t *testing.T) {
	NewArenaWithCapacity[int](10)
}

func TestNewFixedSizeArena(t *testing.T) {
	buffer := [...]int{0, 0, 0}
	metadata := [...]MetadataEntry{{}, {}, {}}

	NewFixedSizeArena[int](buffer[:], metadata[:])
}

func TestArena_Insert(t *testing.T) {
	t.Run("First insert should return zero ID", func(t *testing.T) {
		var a Arena[int]

		id := a.Insert(1)
		expectedID := ID{0, 1}
		if id != expectedID {
			t.Errorf("Insert() returned invalid position")
		}
	})
	t.Run("Second insert should increase position", func(t *testing.T) {
		var a Arena[int]

		a.Insert(1)

		id := a.Insert(2)
		expectedID := ID{1, 1}
		if id != expectedID {
			t.Errorf("Insert() returned invalid position")
		}
	})

	t.Run("Insert after remove should reuse position", func(t *testing.T) {
		var a Arena[int]

		id1 := a.Insert(1)
		a.Insert(2)
		a.Remove(id1)

		id3 := a.Insert(3)
		expectedID := ID{0, 2}
		if id3 != expectedID {
			t.Errorf("Insert() returned invalid position")
		}
	})

	t.Run("Generation should increase after remove", func(t *testing.T) {
		var a Arena[int]

		a.Remove(a.Insert(1))
		a.Remove(a.Insert(2))
		id3 := a.Insert(3)
		expectedID := ID{0, 3}
		if id3 != expectedID {
			t.Errorf("Insert() returned invalid position")
		}
	})
}

func TestArena_Remove(t *testing.T) {
	t.Run("Remove should return false if arena is empty", func(t *testing.T) {
		var a Arena[int]

		if a.Remove(ID{0, 0}) != false {
			t.Errorf("Remove() should fail")
		}
	})

	t.Run("Remove should return true if element is deleted", func(t *testing.T) {
		var a Arena[int]
		id := a.Insert(1)

		if a.Remove(id) != true {
			t.Errorf("Remove() should succeed")
		}
	})

	t.Run("Remove should return false if element is already deleted", func(t *testing.T) {
		var a Arena[int]
		a.Insert(1)
		a.Remove(ID{0, 0})

		if a.Remove(ID{0, 0}) != false {
			t.Errorf("Remove() should fail")
		}
	})
}

func TestArena_Get(t *testing.T) {
	t.Run("Get should return value if element exists", func(t *testing.T) {
		expectedValue := 42

		var a Arena[int]
		a.Insert(1)
		id := a.Insert(expectedValue)

		v, ok := a.Get(id)
		if ok != true {
			t.Errorf("Get() should succeed")
		}

		if v != expectedValue {
			t.Errorf("Get() returned invalid value")
		}
	})

	t.Run("Get should return false if element does not exist", func(t *testing.T) {
		var a Arena[int]

		_, ok := a.Get(ID{0, 0})

		if ok != false {
			t.Errorf("Get() returned invalid value")
		}
	})

	t.Run("Get should return false if element is deleted", func(t *testing.T) {
		var a Arena[int]
		id := a.Insert(42)
		a.Remove(id)

		_, ok := a.Get(id)

		if ok != false {
			t.Errorf("Get() returned invalid value")
		}
	})
}

func TestArena_Len(t *testing.T) {
	t.Run("Len should return 0 if arena is empty", func(t *testing.T) {
		var a Arena[int]
		if a.Len() != 0 {
			t.Errorf("Len() returned invalid length")
		}
	})

	t.Run("Len should return 3 if arena has three elements", func(t *testing.T) {
		var a Arena[int]
		a.Insert(1)
		a.Insert(2)
		a.Insert(3)
		if a.Len() != 3 {
			t.Errorf("Len() returned invalid length")
		}
	})

	t.Run("Len should return 1 if 1 element is removed", func(t *testing.T) {
		var a Arena[int]
		id1 := a.Insert(1)
		a.Insert(2)
		a.Remove(id1)
		if a.Len() != 1 {
			t.Errorf("Len() returned invalid length")
		}
	})
}

func TestArena_Cap(t *testing.T) {
	t.Run("Cap should return 0 if arena is empty", func(t *testing.T) {
		var a Arena[int]
		if a.Cap() != 0 {
			t.Errorf("Cap() returned invalid capacity")
		}
	})

	t.Run("Cap should return 123 if arena has capacity 123", func(t *testing.T) {
		a := NewArenaWithCapacity[int](123)
		if a.Cap() != 123 {
			t.Errorf("Cap() returned invalid capacity")
		}
	})

	t.Run("Cap should increase if arena is growing", func(t *testing.T) {
		a := NewArenaWithCapacity[int](1)
		a.Insert(1)
		a.Insert(2)
		if a.Cap() < 2 {
			t.Errorf("Cap() returned invalid capacity")
		}
	})
}

func TestArena_Clone(t *testing.T) {
	t.Run("Clone should copy all elements", func(t *testing.T) {
		var a Arena[int]
		a.Insert(1)
		a.Insert(2)
		a.Insert(3)

		clone := a.Clone()
		assertValues(t, clone, 1, 2, 3)
	})

	t.Run("Clone should not share memory", func(t *testing.T) {
		var a Arena[int]
		a.Insert(1)
		clone := a.Clone()

		clone.Insert(2)
		assertValues(t, a, 1)
		assertValues(t, clone, 1, 2)
	})
}

func TestArena_Clear(t *testing.T) {
	t.Run("Clear should remove all elements", func(t *testing.T) {
		var a Arena[int]
		a.Insert(1)
		a.Insert(2)
		a.Insert(3)

		a.Clear()
		assertValues(t, a)
	})
}

func TestArena_String(t *testing.T) {
	t.Run("String should return empty string if arena is empty", func(t *testing.T) {
		var a Arena[int]
		expected := "Arena{metadata=[], values=[], freeHead=0}"
		if a.String() != expected {
			t.Errorf("String() returned invalid string")
		}
	})

	t.Run("String should return valid string if arena has elements", func(t *testing.T) {
		var a Arena[int]
		a.Insert(1)
		a.Insert(2)
		a.Insert(3)

		expected := "Arena{metadata=[entry{state=occupied, generation=1} entry{state=occupied, generation=1} entry{state=occupied, generation=1}], values=[1 2 3], freeHead=3}"
		if a.String() != expected {
			t.Errorf("String() returned invalid string, got %v, expected %v", a.String(), expected)
		}
	})
}

func TestScenario(t *testing.T) {
	var a Arena[string]
	id1 := a.Insert("hello")
	id2 := a.Insert("world!")
	id3 := a.Insert("thanks")
	assertValues(t, a, "hello", "world!", "thanks")

	a.Remove(id2)
	id4 := a.Insert("folks!")
	assertValues(t, a, "hello", "folks!", "thanks")

	a.Remove(id1)
	a.Remove(id3)
	assertValues(t, a, "folks!")

	a.Remove(id4)
	assertValues(t, a)
}

func assertValues[T comparable](t *testing.T, a Arena[T], expectedValues ...T) {
	t.Helper()

	if a.Len() != len(expectedValues) {
		t.Errorf("Invalid number of values")
	}

	it := a.Iterate()
	i := 0
	for {
		v, _, done := it.Next()
		if done {
			break
		}

		if v != expectedValues[i] {
			t.Errorf("Iterate returned %v, expected %v", v, expectedValues[i])
		}

		i++
	}
}

func BenchmarkArena_InsertAndRemove(b *testing.B) {
	const n = 1024
	var buffer [n]int
	var metadata [n]MetadataEntry
	var result [n]ID
	a := NewFixedSizeArena[int](buffer[:], metadata[:])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			result[j] = a.Insert(i)
		}
		for j := 0; j < n; j++ {
			a.Remove(result[j])
		}
	}
}
