package madison

import (
	"errors"
	"fmt"
)

// Compute the type of this constant.
func (c Const) Type(cs []CallSite, lcl []Type) (Type, error) {
	return Type{Range: Range{int(c), int(c)}}, nil
}

// Attempt to set the type of this constant.
func (c Const) RestrictTo(locals []Type, t Type) error {
	if t.Elem != nil || t.Range.Start > int(c) || t.Range.End < int(c) {
		return &Impossible{c, t, Type{Range{int(c), int(c)}, nil}}
	} else {
		return nil
	}
}

// Compute the type of the empty list.
func (e EmptyList) Type(cs []CallSite, lcl []Type) (Type, error) {
	return Type{Range: Range{0, 0}, Elem: &Type{Range: UNDEF}}, nil
}

// Attempt to set the type of the empty list.
func (e EmptyList) RestrictTo(locals []Type, t Type) error {
	if t.Elem == nil || t.Elem.Range.Start != 0 || t.Elem.Range.End != 0 {
		return &Impossible{e, t, Type{Range{0, 0}, &Type{Range: UNDEF}}}
	} else {
		return nil
	}
}

// Compute the type of the sum of the two arguments.
func (p *Plus) Type(cs []CallSite, lcl []Type) (Type, error) {
	a, err := p.A.Type(cs, lcl)
	if err != nil {
		return a, err
	}
	b, err := p.B.Type(cs, lcl)
	if err != nil {
		return b, err
	}
	if a.Elem != nil {
		fmt.Printf("can't add lists")
	}
	return Type{Range: conv(a.Range, b.Range)}, nil
}

// Attempt to set the type of the sum of the arguments.
func (p *Plus) RestrictTo(locals []Type, t Type) error {
	// Try to set A to T - B
	b, err := p.B.Type([]CallSite{}, locals)
	if err != nil {
		return err
	}
	aerr := p.A.RestrictTo(locals,
		Type{Range: conv(t.Range, inverse(b.Range))})

	// Try to set B to T - A
	a, err := p.A.Type([]CallSite{}, locals)
	if err != nil {
		return err
	}
	berr := p.B.RestrictTo(locals,
		Type{Range: conv(t.Range, inverse(a.Range))})

	if aerr != nil {
		return berr
	} else {
		return aerr
	}
}

// Compute the type of this negation.
func (n *Negate) Type(cs []CallSite, lcl []Type) (Type, error) {
	typ, err := n.Elem.Type(cs, lcl)
	if err != nil {
		return NIL, err
	}
	if typ.Elem != nil {
		return NIL, fmt.Errorf("cannot negate a %s", typ)
	}
	return Type{
		Range: inverse(typ.Range),
	}, nil
}

// Attempt to set the type of the negation.
func (n *Negate) RestrictTo(locals []Type, t Type) error {
	if t.Elem != nil {
		return errors.New("negating an array is not supported yet")
	}

	t.Range = inverse(t.Range)
	return n.Elem.RestrictTo(locals, t)
}

// Compute the type of this variable reference.
func (v *Var) Type(cs []CallSite, lcl []Type) (Type, error) {
	return lcl[v.index], nil
}

// Attempt to set the type of this variable.
func (v *Var) RestrictTo(locals []Type, t Type) error {
	intr := intersect(locals[v.index].Range, t.Range)
	if len(intr) == 0 {
		return &Impossible{v, t, locals[v.index]}
	}
	locals[v.index].Range = intr[0]
	return nil
}

// Compute the type of this conditional.
func (i *If) Type(cs []CallSite, lcl []Type) (Type, error) {
	var (
		lte Type
		gte Type
	)

	// <= 0 case:
	copy := append([]Type(nil), lcl...)
	lerr := i.Cond.RestrictTo(copy, NON_POSITIVE)
	if lerr == nil {
		if lte, lerr = i.NonPositive.Type(cs, copy); lerr != nil {
			return NIL, lerr
		}
	}

	// >0 case:
	copy = append([]Type(nil), lcl...)
	gerr := i.Cond.RestrictTo(copy, POSITIVE)
	if gerr == nil {
		if gte, gerr = i.Positive.Type(cs, copy); gerr != nil {
			return NIL, gerr
		}
	}

	if lerr == nil && gerr == nil {
		return TypesUnion(lte, gte)
	} else if lerr == nil {
		return lte, nil
	} else {
		return gte, gerr
	}
}

// Attempt to set the type of this conditional.
func (i *If) RestrictTo(locals []Type, t Type) error {
	copy := append([]Type(nil), locals...)
	if err := i.Cond.RestrictTo(copy, NON_POSITIVE); err == nil {
		if err := i.NonPositive.RestrictTo(copy, t); err == nil {
			if err := i.Cond.RestrictTo(locals, NON_POSITIVE); err != nil {
				return err
			}
		}
	}

	copy = append([]Type(nil), locals...)
	if err := i.Cond.RestrictTo(copy, POSITIVE); err == nil {
		if err := i.Positive.RestrictTo(copy, t); err == nil {
			if err := i.Cond.RestrictTo(locals, POSITIVE); err != nil {
				return err
			}
		}
	}

	return nil
}

// Compute the type of this function call.
func (a *Apply) Type(cs []CallSite, locals []Type) (Type, error) {
	arg, err := a.Arg.Type(cs, locals)
	if err != nil {
		return NIL, err
	}

	funct, ok := a.Runtime.Funcs[a.Name]
	if !ok {
		return NIL, fmt.Errorf("undefined function %s", a.Name)
	}

	return funct.Type(cs, []Type{arg})
}

// Attempt to set the type of this function call.
func (a *Apply) RestrictTo(locals []Type, t Type) error {
	arg, err := a.Arg.Type([]CallSite{}, locals)
	if err != nil {
		return err
	}

	funct, ok := a.Runtime.Funcs[a.Name]
	if !ok {
		return fmt.Errorf("undefined function %#v\n", a.Name)
	}

	lcls := []Type{arg}
	if err := funct.RestrictTo(lcls, t); err != nil {
		return err
	}
	return a.Arg.RestrictTo(locals, lcls[0])
}

// Compute the type of this prepend call.
func (p *Prepend) Type(cs []CallSite, locals []Type) (Type, error) {
	h, err := p.Head.Type(cs, locals)
	if err != nil {
		return NIL, err
	}
	t, err := p.Tail.Type(cs, locals)
	if err != nil {
		return NIL, err
	}
	if t.Elem == nil {
		return NIL, fmt.Errorf("%s is not a list type", t)
	}
	m := h
	if t.Range.End > 0 {
		if m, err = TypesUnion(h, *t.Elem); err != nil {
			return NIL, err
		}
	}
	return Type{
		Range: conv(Range{1, 1}, t.Range),
		Elem:  &m,
	}, nil
}

// Attempt to set the type of this prepend call.
func (p *Prepend) RestrictTo(locals []Type, t Type) error {
	if t.Elem != nil || t.Elem.Start < 1 {
		return &Impossible{p, t, Type{Range{1, 1}, &Type{}}}
	}
	if err := p.Head.RestrictTo(locals, *t.Elem); err != nil {
		return err
	}
	t.Range.Start -= 1
	return p.Tail.RestrictTo(locals, t)
}

// Compute the type of the first element of the list.
func (h *Head) Type(cs []CallSite, locals []Type) (Type, error) {
	typ, err := h.List.Type(cs, locals)
	if err != nil {
		return NIL, err
	} else if typ.Elem == nil {
		return NIL, fmt.Errorf("element is not a list type: %s", typ)
	} else if typ.Range.Start < 1 {
		return NIL, fmt.Errorf("cannot take head of an empty list: %s", typ)
	}
	return *typ.Elem, nil
}

// Attempt to set the type of the first element of this list.
func (h *Head) RestrictTo(locals []Type, t Type) error {
	return h.List.RestrictTo(locals, Type{Range: POSITIVE.Range, Elem: &t})
}

// Compute the type of the remaining elements of the list.
func (t *Tail) Type(cs []CallSite, locals []Type) (Type, error) {
	typ, err := t.List.Type(cs, locals)
	if err != nil {
		return NIL, err
	} else if typ.Elem != nil {
		return NIL, fmt.Errorf("element is not a list type: %s", typ)
	} else if typ.Range.Start < 1 {
		return NIL, fmt.Errorf("cannot take tail of an empty list: %s", typ)
	}
	return *typ.Elem, nil
}

// Attempt to set the type of the remaining elements of the list.
func (t *Tail) RestrictTo(locals []Type, typ Type) error {
	if typ.Elem != nil {
		return fmt.Errorf("element is not a list type: %s", typ)
	}
	typ.Range.Start -= 1
	return t.List.RestrictTo(locals, typ)
}

// Raises a pattern match failure.
func (t *Undef) Type(cs []CallSite, locals []Type) (Type, error) {
	return NIL, errors.New("undefined")
}

// Raises a pattern match failure.
func (t *Undef) RestrictTo(locals []Type, typ Type) error {
	return errors.New("undefined")
}
