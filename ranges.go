package madison

import (
	"fmt"
	"math"
)

// A range of integers.
type Range struct {
	Start, End int
}

// Represents any possible integer.
var UNDEF = Range{math.MinInt32, math.MaxInt32}

// Pretty-prints this Range.
func (r Range) String() string {
	if r.Start == math.MinInt32 && r.End == math.MaxInt32 {
		return "any"
	}
	if r.Start == r.End {
		return fmt.Sprintf("%d", r.Start)
	}
	st := "(-∞"
	if r.Start > math.MinInt32 {
		st = fmt.Sprintf("[%d", r.Start)
	}
	ed := "∞)"
	if r.End < math.MaxInt32 {
		ed = fmt.Sprintf("%d]", r.End)
	}
	return fmt.Sprintf("%s, %s", st, ed)
}

func inverse(r Range) Range {
	a, b := -r.End, -r.Start
	if r.Start == math.MinInt32 {
		b = math.MaxInt32
	}
	if r.End == math.MaxInt32 {
		a = math.MinInt32
	}
	return Range{a, b}
}

func union(a, b Range) (o Range) {
	o = b
	if a.Start < b.Start {
		o.Start = a.Start
	}
	if a.End > b.End {
		o.End = a.End
	}
	return
}

func conv(a, b Range) Range {
	st := a.Start + b.Start
	if a.Start == math.MinInt32 || b.Start == math.MinInt32 {
		st = math.MinInt32
	}
	ed := a.End + b.End
	if a.End == math.MaxInt32 || b.End == math.MaxInt32 {
		ed = math.MaxInt32
	}
	return Range{st, ed}
}

func (r Range) IsConst() bool {
	return r.Start == r.End // Start cannot be +inf, and End cannot be -inf
}

func intersect(a Range, b Range) []Range {
	if a.End < b.Start || b.End < a.Start {
		return []Range{}
	}
	if b.Start > a.Start {
		a.Start = b.Start
	}
	if b.End < a.End {
		a.End = b.End
	}
	return []Range{a}
}

func subtract(a Range, b Range) []Range {
	if a.Start < b.Start && b.End < a.End {
		return []Range{Range{a.Start, b.Start - 1}, Range{b.End + 1, a.End}}
	}
	if b.Start <= a.Start && a.End <= b.End {
		return []Range{}
	}
	if a.End < b.Start || b.End < a.Start {
		return []Range{a} // no overlap
	}
	if a.Start < b.Start {
		a.End = b.Start - 1
	} else {
		a.Start = b.End + 1
	}
	if a.Start > a.End {
		panic(fmt.Sprintf("creating invalid range: %s", a))
	}
	return []Range{a}
}
