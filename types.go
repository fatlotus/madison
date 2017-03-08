package madison

import (
	"fmt"
	"math"
)

// Represent a callsite (not implemented yet).
type CallSite struct{}

// Represents a type.
type Type struct {
	// How large this value can get.
	Range

	// If non-nil: this Type is a list containing Elem elemnts.
	// else: this type is a scalar.
	Elem *Type
}

// Represents a type mismatch.
type Impossible struct {
	// Which node failed the pattern match.
	Context Node

	// What the given node evaluated to.
	Found Type

	// What we needed the node to evaluate to.
	Needed Type
}

// Represent the Type error as an error.
func (i *Impossible) Error() string {
	return fmt.Sprintf("needed %s, but got %s, in %s",
		i.Found, i.Needed, i.Context)
}

var (
	// Represents x <= 0.
	NON_POSITIVE = InRange(math.MinInt32, 0)

	// Represents x > 0.
	POSITIVE = InRange(1, math.MaxInt32)

	// Represents x = 0.
	NIL = Constant(0)
)

// Return a scalar type with the single constant value.
func Constant(v int) Type {
	return Type{Range: Range{v, v}}
}

// Return a scalar type that exists in the given range.
func InRange(min, max int) Type {
	return Type{Range: Range{min, max}}
}

// Pretty-prints a type.
//
// Representation:
//               v precision
//   [2, 5][1, 2]int[3, 4]
//   ^ dimensions   ^ range of values
//
func (t Type) String() string {
	if t.Elem != nil {
		r := t.Range.String()
		if r[0] != '[' && r[0] != '(' {
			r = "[" + r + "]"
		}
		return r + t.Elem.String()
	}
	s := t.Range.String()
	if s[0] == '[' || s[0] == '(' {
		s = "int" + s
	}
	return s
}

// Returns true if the given type is a subset of another.
func (t Type) SubsetOf(o Type) bool {
	if !(t.Start <= o.End && o.Start <= t.End) {
		return false
	}
	return (t.Start < o.Start || o.End < t.End ||
		(t.Elem != nil && o.Elem != nil && t.Elem.SubsetOf(*o.Elem)))
}

// Joins the two types together (making one that is less specific than either).
func TypesUnion(a, b Type) (Type, error) {
	t := Type{Range: union(a.Range, b.Range)}
	if a.Elem != nil {
		if a.Range.End == 0 {
			t.Elem = b.Elem
		} else if b.Range.End == 0 {
			t.Elem = a.Elem
		} else {
			u, _ := TypesUnion(*a.Elem, *b.Elem)
			t.Elem = &u
		}
	}
	return t, nil // add loads more checking
}
