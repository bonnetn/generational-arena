package arena

import "fmt"

// state represents the state of a metadata entry.
type state bool

const (
	// stateFree represents a free metadata entry.
	stateFree     state = false
	// stateOccupied represents an occupied metadata entry.
	stateOccupied state = true
)

// MetadataEntry represents an entry in the metadata of the arena.
type MetadataEntry struct {
	state      state
	generation int64
	nextHead   int // Populated if free
}

// String returns a string representation of the metadata entry.
func (e MetadataEntry) String() string {
	if e.state == stateOccupied {
		return fmt.Sprintf("entry{state=occupied, generation=%d}", e.generation)
	} else {
		return fmt.Sprintf("entry{state=free, nextHead=%d}", e.nextHead)
	}
}
