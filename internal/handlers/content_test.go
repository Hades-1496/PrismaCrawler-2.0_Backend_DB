package handlers

import (
	"strings"
	"testing"
)

func TestGenerateMapInvariants(t *testing.T) {
	seeds := []string{"", "a", "zzzz", "12345", "seed-x"}
	for _, s := range seeds {
		for _, lvl := range []int{1, 5, 12, 20} {
			layout, _ := generateMap(lvl, s)
			joined := strings.Join(layout, "")
			if strings.Count(joined, "P") != 1 {
				t.Fatalf("nivel %d seed %q: se esperaba 1 'P', hubo %d", lvl, s, strings.Count(joined, "P"))
			}
			if strings.Count(joined, "E") < 1 {
				t.Fatalf("nivel %d seed %q: sin salida 'E'", lvl, s)
			}
		}
	}
}
