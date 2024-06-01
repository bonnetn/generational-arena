package arena

import (
	"fmt"
	"testing"
)

func TestIdToMnemomic(t *testing.T) {
	var a Arena[string]
	for i := 0; i < 100; i++ {
		id := a.Insert("hello")
		m := id.Mnemomic()
		a.Remove(id)

		fmt.Println(m)
		fmt.Println(id)
	}
}
func TestID_Mnemomic(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			name: "(0,1)",
			id:   ID{Position: 0, Generation: 1},
			want: "acrobat",
		},
		{
			name: "(1,1)",
			id:   ID{Position: 1, Generation: 1},
			want: "carbon",
		},
		{
			name: "(2,1)",
			id:   ID{Position: 2, Generation: 1},
			want: "humor",
		},
		{
			name: "(3,1)",
			id:   ID{Position: 3, Generation: 1},
			want: "total",
		},
		{
			name: "(1,2)",
			id:   ID{Position: 1, Generation: 2},
			want: "casino",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.Mnemomic(); got != tt.want {
				t.Errorf("ID.Mnemomic() = %v, want %v", got, tt.want)
			}
		})
	}
}
