package arena

import "testing"

func TestID_String(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		i := ID{Position: 1, Generation: 1}
		expected := "ID{index=1, generation=1}"
		if i.String() != expected {
			t.Errorf("String() returned invalid string")
		}
	})
}
