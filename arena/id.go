package arena

import "fmt"

// ID represents an identifier for an element in the arena.
type ID struct {
	Position   int
	Generation int64
}

// String returns a string representation of the ID.
func (i ID) String() string {
	return fmt.Sprintf("ID{index=%d, generation=%d}", i.Position, i.Generation)
}
