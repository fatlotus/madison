package madison

import (
	"fmt"
)

// Represents a node in the tree (i.e. a thing that, if it has a type, can be
// evaluated).
type Node interface {
	// Evaluates this function with the given local variables.
	Eval(lcl []Obj) Obj

	// Computes the type of this Node given the local arguments.
	Type(callers []CallSite, locals []Type) (Type, error)

	// Attempts to update the given arguments so that the result is a subset
	// of the given type.
	RestrictTo(locals []Type, t Type) error

	String() string
}

// Constant values.
type Const int

var _ Node = Const(0)

// Prints this constant.
func (c Const) String() string { return fmt.Sprintf("%d", c) }

// An empty list (i.e. [0]any).
type EmptyList struct{}

var _ Node = EmptyList{}

// Prints the empty list.
func (n EmptyList) String() string { return "[]" }

// binary math ops (a + b)
type Plus struct{ A, B Node }

var _ Node = &Plus{}

// Prints the sum of the two arguments.
func (p *Plus) String() string {
	if neg, ok := p.B.(*Negate); ok {
		return fmt.Sprintf("(%s - %s)", p.A, neg.Elem)
	}
	return fmt.Sprintf("(%s + %s)", p.A, p.B)
}

// Negates the given node.
type Negate struct{ Elem Node }

var _ Node = &Negate{}

// Prints the opposite of Elem.
func (n *Negate) String() string {
	return fmt.Sprintf("-%s", n.Elem)
}

// referring to a function argument
type Var struct{ index int }

const vars = "xyzwabc"

var _ Node = &Var{}

// Pretty-prints this variable.
func (v *Var) String() string {
	if v.index < len(vars) {
		return vars[v.index : v.index+1]
	}
	return fmt.Sprintf("@%d", v.index)
}

// A conditional block.
type If struct{ Cond, NonPositive, Positive Node }

var _ Node = &If{}

// Pretty-prints this conditional block.
func (n *If) String() string {
	return fmt.Sprintf("ifz(%s, %s, %s)", n.Cond, n.NonPositive, n.Positive)
}

// Puts the given Head on the fron of the Tail list.
type Prepend struct {
	Head, Tail Node
}

var _ Node = &Prepend{}

// Pretty-prints this Prepend cell.
func (p *Prepend) String() string {
	return fmt.Sprintf("%s : %s", p.Head, p.Tail)
}

// Computes the first element in List.
type Head struct {
	List Node
}

var _ Node = &Head{}

// Pretty-prints this Head cell.
func (h *Head) String() string {
	return fmt.Sprintf("head(%s)", h.List)
}

// Represents nodes in List after the first.
type Tail struct {
	List Node
}

var _ Node = &Tail{}

// Pretty-prints this Tail.
func (h *Tail) String() string {
	return fmt.Sprintf("tail(%s)", h.List)
}

// Represents unconditional failure (represents pattern match failure).
type Undef struct {
	Message string
}

var _ Node = &Undef{}

// Represents pattern match failure.
func (u *Undef) String() string {
	return "undef()"
}

// Stores all named functions in the runtime.
type Runtime struct {
	Funcs map[string]Node
}

// Call a function in the current runtime by name.
type Apply struct {
	Runtime *Runtime
	Name    string
	Arg     Node
}

var _ Node = &Apply{}

// Pretty-print this Apply call.
func (a *Apply) String() string {
	return fmt.Sprintf("%s(%s)", a.Name, a.Arg)
}
