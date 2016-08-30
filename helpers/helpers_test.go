package helpers

import "testing"

func TestArithmetic(t *testing.T) {
	t.Run("Mul", func(t *testing.T) {
		if got := Mul(12, 3); got != 36 {
			t.Errorf("Got Mul(12,3)=%d, want 36", got)
		}
	})

	t.Run("Add", func(t *testing.T) {
		if got := Add(100, 50); got != 150 {
			t.Errorf("Got Mul(100,50)=%d, want 150", got)
		}
	})
}

func TestMakeSlace(t *testing.T) {

	testCases := [...]interface{}{
		1,
		2,
		"3",
		nil,
		struct{}{},
	}
	got := MkSlice(testCases[:]...)
	if len(got) != len(testCases) {
		t.Fatalf("Got MkSlice(%v)=%v,want %v", got, testCases, got)
	}
	for i, w := range testCases {
		if w != got[i] {
			t.Errorf("Got MkSlice(%v)=%v,want %v", got, testCases, got)
		}
	}

	if got := MkSlice(nil); len(got) != 1 || got[0] != nil {
		t.Errorf("Got MkSlice(%v)=%v,want len = 1 and [%v] ", nil, got, nil)
	}

}

func IntsEquals(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestMakeSlaceRange(t *testing.T) {
	testCases := []struct {
		infrom  int
		incount int
		out     []int
	}{
		{0, 2, []int{0, 1}},
		{5, 5, []int{5, 6, 7, 8, 9}},
		{1, 1, []int{1}},
		{5, 0, nil},
	}

	for _, tc := range testCases {
		if got := MkSliceRange(tc.infrom, tc.incount); !IntsEquals(got, tc.out) {
			t.Errorf("Got MkSliceRange(%d,%d)=%v, want %v", tc.infrom, tc.incount, got, tc)
		}
	}
}

func TestCalcPages(t *testing.T) {
	testCases := []struct{ onpage, count, want int }{
		{65, 66, 2},
		{65, 10, 1},
		{65, 150, 3},
		{65, 150, 3},
		{0, 0, 0},
	}
	for _, tc := range testCases {
		if got := Calcpages(tc.onpage, tc.count); got != tc.want {
			t.Errorf("Got Calcpages(%d,%d)=%v, want %v", tc.onpage, tc.count, got, tc.want)
		}
	}

}
