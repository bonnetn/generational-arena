package arena

// Iterator is an Iterate over the values in an arena.
type Iterator[T any] struct {
	i     int
	arena *Arena[T]
}

// Next returns the next value in the Iterate, along with its ID.
// If there are no more values, it returns true for done.
func (e *Iterator[T]) Next() (val T, id ID, done bool) {
	for e.i < len(e.arena.metadata) {
		if e.arena.metadata[e.i].state == stateOccupied {
			v := e.arena.values[e.i]
			id := ID{e.i, e.arena.metadata[e.i].generation}
			e.i++
			return v, id, false
		}
		e.i++
	}

	return val, id, true
}
