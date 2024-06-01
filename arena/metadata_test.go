package arena

import (
	"testing"
)

func TestMetadata_String(t *testing.T) {
	t.Run("occupied", func(t *testing.T) {
		m := MetadataEntry{state: stateOccupied, generation: 1}
		expected := "entry{state=occupied, generation=1}"
		if m.String() != expected {
			t.Errorf("String() returned invalid string")
		}
	})

	t.Run("free", func(t *testing.T) {
		m := MetadataEntry{state: stateFree, nextHead: 1}
		expected := "entry{state=free, nextHead=1}"
		if m.String() != expected {
			t.Errorf("String() returned invalid string")
		}
	})
}
