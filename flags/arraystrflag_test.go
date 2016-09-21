package flags

import "testing"
import "fmt"

func TestArrayString(t *testing.T) {
	var arr ArrayString
	if got := arr.Count(); got != 0 {
		t.Fatalf("Got Count()=%d, want %d", got, 0)
	}
	err := arr.Set("test")
	if err != nil {
		t.Fatal(err)
	}
	if got := arr.Count(); got != 1 {
		t.Fatalf("Got Count()=%d, want %d", got, 1)
	}

	if got := arr.String(); got != fmt.Sprintf("%s", []string(arr)) {
		t.Errorf("Got String()=%s,want %s", got, fmt.Sprintf("%s", []string(arr)))
	}
}
