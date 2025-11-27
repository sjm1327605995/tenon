package yoga

import "testing"

func TestInexactEquals(t *testing.T) {
	tests := []struct {
		a, b   any
		expect bool
	}{
		{float32(1), float32(1), true},
		{float32(1), float32(2), false},
		{float32(1), float64(1.00011), false},
		{float64(1), float64(1), true},
		{float64(1), float64(2), false},
		{float64(1), float64(1.0001), true},
		{float64(1), float64(1.00011), false},
		{Value{1, UnitPoint}, Value{1, UnitPoint}, true},
		{Value{1, UnitPoint}, Value{2, UnitPoint}, false},
		{Value{1, UnitPoint}, Value{1.00011, UnitPoint}, false},

		{FloatOptional{1}, FloatOptional{1}, true},
		{FloatOptional{1}, FloatOptional{2}, false},
		{FloatOptional{1}, FloatOptional{1.00011}, false},
	}
	for _, test := range tests {
		if inexactEquals(test.a, test.b) != test.expect {
			t.Errorf("inexactEquals(%v, %v) != %v", test.a, test.b, test.expect)
		}
	}
}
