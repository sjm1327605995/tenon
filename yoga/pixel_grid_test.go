package yoga

import "testing"

func TestRoundValueToPixelGrid(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		scale       float64
		ceil, floor bool
		want        float32
	}{
		{name: "nearest down", value: 10.25, scale: 1, want: 10},
		{name: "nearest up", value: 10.75, scale: 1, want: 11},
		{name: "negative nearest", value: -2.2, scale: 1, want: -2},
		{name: "negative floor", value: -2.2, scale: 1, floor: true, want: -3},
		{name: "negative ceil", value: -2.2, scale: 1, ceil: true, want: -2},
		{name: "scaled", value: 10.25, scale: 2, want: 10.5},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := RoundValueToPixelGrid(test.value, test.scale, test.ceil, test.floor)
			if got != test.want {
				t.Fatalf("RoundValueToPixelGrid()=%v, want %v", got, test.want)
			}
		})
	}
}
