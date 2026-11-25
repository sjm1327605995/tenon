package g3n

import "testing"

func TestSmokeFrameLimit(t *testing.T) {
	for _, test := range []struct {
		value string
		want  int
	}{
		{value: "", want: 0},
		{value: "invalid", want: 0},
		{value: "0", want: 0},
		{value: "3", want: 3},
	} {
		t.Run(test.value, func(t *testing.T) {
			t.Setenv("TENON_SMOKE_FRAMES", test.value)
			if got := smokeFrameLimit(); got != test.want {
				t.Fatalf("smokeFrameLimit()=%d, want %d", got, test.want)
			}
		})
	}
}
