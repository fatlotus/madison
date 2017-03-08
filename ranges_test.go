package madison

import (
	"math"
	"testing"
)

var (
	ninf = math.MinInt32
	inf  = math.MaxInt32
)

var singles = []struct {
	Op  func(r Range) Range
	Arg Range
	Exp Range
}{
	{inverse, Range{1, 2}, Range{-2, -1}},
	{inverse, Range{1, inf}, Range{ninf, -1}},
	{inverse, Range{ninf, inf}, Range{ninf, inf}},
	{inverse, Range{ninf, 3}, Range{-3, inf}},
}

var doubles = []struct {
	Op   func(a Range, b Range) Range
	Arg1 Range
	Arg2 Range
	Exp  Range
}{
	{union, Range{1, 3}, Range{-1, 2}, Range{-1, 3}},
	{union, Range{1, inf}, Range{-1, 2}, Range{-1, inf}},
	{union, Range{ninf, 3}, Range{-1, 2}, Range{ninf, 3}},
	{union, Range{ninf, 3}, Range{-1, inf}, Range{ninf, inf}},

	{conv, Range{1, 2}, Range{1, 1}, Range{2, 3}},
	{conv, Range{1, 2}, Range{1, 2}, Range{2, 4}},
	{conv, Range{ninf, 1}, Range{1, 2}, Range{ninf, 3}},
	{conv, Range{ninf, 1}, Range{1, inf}, Range{ninf, inf}},
	{conv, Range{-1, 1}, Range{1, inf}, Range{0, inf}},
}

func TestRange(t *testing.T) {
	for i, c := range singles {
		got := c.Op(c.Arg)
		if got != c.Exp {
			t.Errorf("%d: (%s) = %s (expecting %s)", i, c.Arg, got, c.Exp)
			t.Fail()
		}
	}
	for i, c := range doubles {
		got := c.Op(c.Arg1, c.Arg2)
		if got != c.Exp {
			t.Errorf("%d: (%s, %s) = %s (expecting %s)", i, c.Arg1, c.Arg2, got, c.Exp)
		}
	}
}
