package arena

import (
	"fmt"
)

// Arena is a generational arena that stores values of a given type.
// The zero value of an arena is an empty arena.
// The arena is not thread-safe.
type Arena[T any] struct {
	metadata []MetadataEntry
	values   []T
	freeHead int
	length   int
	size     int
}

// NewArenaWithCapacity creates a new arena with a given capacity.
func NewArenaWithCapacity[T any](capacity int) Arena[T] {
	return Arena[T]{
		metadata: make([]MetadataEntry, 0, capacity),
		values:   make([]T, 0, capacity),
	}
}

// NewFixedSizeArena creates a new arena using the slices provided as input.
// The length of the slices must be equal and strictly greater than 0.
// The resulting arena will panic if the capacity is reached.
// This allows for zero allocation arenas.
func NewFixedSizeArena[T any](buffer []T, metadata []MetadataEntry) Arena[T] {
	if len(buffer) != len(metadata) {
		panic("buffer and metadata length must be equal")
	}

	if len(buffer) == 0 {
		panic("buffer and metadata length must be greater than 0")
	}

	return Arena[T]{
		values:   buffer[:0],
		metadata: metadata[:0],
		size:     len(buffer),
	}
}

// Insert adds a value to the arena and returns an ID that can be used to retrieve or remove the value.
func (a *Arena[T]) Insert(v T) ID {
	if a.freeHead < len(a.metadata) {
		insertPosition := a.freeHead
		currentEntry := &a.metadata[insertPosition]

		currentEntry.state = stateOccupied
		currentEntry.generation += 1
		a.values[insertPosition] = v
		a.freeHead = currentEntry.nextHead
		a.length += 1

		return ID{insertPosition, currentEntry.generation}
	} else {
		if a.size > 0 && len(a.metadata) >= a.size {
			panic("arena is full")
		}

		newEntry := MetadataEntry{state: stateOccupied, generation: 1}
		a.values = append(a.values, v)
		a.metadata = append(a.metadata, newEntry)
		a.freeHead += 1
		a.length += 1

		return ID{len(a.metadata) - 1, newEntry.generation}
	}
}

// Remove deletes a value from the arena.
// It returns true if the value was removed, false otherwise.
func (a *Arena[T]) Remove(id ID) bool {
	if id.Position < 0 || id.Position >= len(a.metadata) {
		return false
	}

	currentEntry := &a.metadata[id.Position]
	if currentEntry.state == stateFree || currentEntry.generation != id.Generation {
		return false
	}

	currentEntry.state = stateFree
	currentEntry.nextHead = a.freeHead
	a.freeHead = id.Position
	a.length -= 1

	return true
}

// Get returns the value associated with the given ID.
// It returns the value and true if the value exists, false otherwise.
func (a *Arena[T]) Get(id ID) (result T, ok bool) {
	if id.Position < 0 || id.Position >= len(a.metadata) {
		return result, false
	}

	existingEntry := a.metadata[id.Position]
	if existingEntry.state == stateFree || existingEntry.generation != id.Generation {
		return result, false
	}

	return a.values[id.Position], true
}

func (a *Arena[T]) Put(id ID, value T) bool {
	if id.Position < 0 || id.Position >= len(a.metadata) {
		return false
	}

	existingEntry := a.metadata[id.Position]
	if existingEntry.state == stateFree || existingEntry.generation != id.Generation {
		return false
	}

	a.values[id.Position] = value
	return true
}

// Clone creates a deep copy of the arena.
// The new arena will have the same capacity as the original arena.
func (a *Arena[T]) Clone() Arena[T] {
	clone := NewArenaWithCapacity[T](cap(a.metadata))
	clone.metadata = append(clone.metadata, a.metadata...)
	clone.values = append(clone.values, a.values...)
	clone.freeHead = a.freeHead
	clone.length = a.length
	clone.size = a.size

	return clone
}

// Len returns the number of elements in the arena.
func (a *Arena[T]) Len() int {
	return a.length
}

// Cap returns the capacity of the arena.
func (a *Arena[T]) Cap() int {
	return cap(a.metadata)
}

// Clear removes all elements from the arena.
// The capacity of the arena remains the same.
func (a *Arena[T]) Clear() {
	a.metadata = a.metadata[:0]
	a.values = a.values[:0]
	a.freeHead = 0
	a.length = 0
}

// Iterate returns an Iterate for the arena.
func (a *Arena[T]) Iterate() *Iterator[T] {
	return &Iterator[T]{arena: a}
}

// String returns a string representation of the arena.
func (a *Arena[T]) String() string {
	return fmt.Sprintf("Arena{metadata=%v, values=%v, freeHead=%d}", a.metadata, a.values, a.freeHead)
}
