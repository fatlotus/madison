package madison

import (
	"fmt"
)

type Obj struct {
	Int  int64
	Vals []Obj
}

func (o Obj) String() string {
	if o.Vals == nil {
		return fmt.Sprintf("%d", o.Int)
	} else {
		buf := ""
		for _, elem := range o.Vals {
			if buf != "" {
				buf += " : "
			}
			buf += elem.String()
		}
		buf += " : []"
		return buf
	}
}

// Evaluate this constant.
func (c Const) Eval(args []Obj) Obj {
	return Obj{Int: int64(int(c))}
}

// Evaluate the empty list.
func (e EmptyList) Eval(args []Obj) Obj {
	return Obj{Vals: []Obj{}}
}

// Evaluate the sum of the two arguments.
func (p *Plus) Eval(args []Obj) Obj {
	return Obj{Int: p.A.Eval(args).Int + p.B.Eval(args).Int}
}

// Evaluate this negation.
func (n *Negate) Eval(args []Obj) Obj {
	return Obj{Int: -n.Elem.Eval(args).Int}
}

// Evaluate this variable reference.
func (v *Var) Eval(args []Obj) Obj {
	return args[v.index]
}

// Evaluate this conditional.
func (i *If) Eval(args []Obj) Obj {
	if i.Cond.Eval(args).Int <= 0 {
		return i.NonPositive.Eval(args)
	} else {
		return i.Positive.Eval(args)
	}
}

// Evaluate this function call.
func (a *Apply) Eval(args []Obj) Obj {
	return a.Runtime.Funcs[a.Name].Eval([]Obj{a.Arg.Eval(args)})
}

// Evaluate this prepend call.
func (p *Prepend) Eval(args []Obj) Obj {
	return Obj{
		Vals: append([]Obj{p.Head.Eval(args)}, p.Tail.Eval(args).Vals...),
	}
}

// Evaluate the first element of the list.
func (h *Head) Eval(args []Obj) Obj {
	return h.List.Eval(args).Vals[0]
}

// Evaluate the remaining elements of the list.
func (t *Tail) Eval(args []Obj) Obj {
	return Obj{Vals: t.List.Eval(args).Vals[1:]}
}

// Evaluate a pattern match failure.
func (t *Undef) Eval(args []Obj) Obj {
	panic("pattern match!")
}
