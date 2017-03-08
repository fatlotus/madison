package madison

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/fatlotus/mast"
)

var parser = mast.Parser{
	Parens: []mast.Group{
		{"(", ")"},
	},
	Brackets: []mast.Group{
		{"[", "]"},
	},
	Operators: []mast.Prec{
		{[]string{","}, mast.InfixRight},
		{[]string{":"}, mast.InfixRight},
		{[]string{"+", "-"}, mast.InfixLeft},
		{[]string{"*"}, mast.InfixLeft},
	},
	AdjacentIsApplication: true,
}

// b 0 = 2 => b = (case @0 of 2 => | a => nil)
func (r *Runtime) mastToTuple(e mast.Expr, lv bool, as *[]string) (a []Node) {
	for {
		t, ok := e.(*mast.Binary)
		if ok && t.Op == "," {
			a = append(a, r.mastToExpr(t.Left, lv, as))
			e = t.Right
		} else {
			break
		}
	}
	return append(a, r.mastToExpr(e, lv, as))
}

func (r *Runtime) mastToExpr(e mast.Expr, lval bool, args *[]string) Node {
	switch e := e.(type) {
	case *mast.Unary: // only -
		x := r.mastToExpr(e.Elem, lval, args)
		return &Negate{x}
	case *mast.Binary:
		if e.Op == ":" { // prepend / cons
			head := r.mastToExpr(e.Left, lval, args)
			tail := r.mastToExpr(e.Right, lval, args)
			return &Prepend{head, tail}
		} else { // + or -
			a := r.mastToExpr(e.Left, lval, args)
			b := r.mastToExpr(e.Right, lval, args)
			if e.Op == "-" {
				b = &Negate{b}
			}
			return &Plus{a, b}
		}
	case *mast.Apply:
		m := e.Operator.(*mast.Var)
		args := r.mastToTuple(e.Operand, lval, args)
		switch m.Name {
		case "ifz":
			if len(args) != 3 {
				panic(fmt.Sprintf("ifz takes three arguments, got %#v", args))
			}
			cond, zero, nonzero := args[0], args[1], args[2]
			return &If{cond, zero, nonzero}
		case "head":
			if len(args) != 1 {
				panic(fmt.Sprintf("head takes one arguments, got %#v", args))
			}
			return &Head{args[0]}
		case "tail":
			if len(args) != 1 {
				panic(fmt.Sprintf("tail takes one arguments, got %#v", args))
			}
			return &Tail{args[0]}
		default:
			if len(args) != 1 {
				panic(fmt.Sprintf(
					"user-defined function %s can only accept 1 arg, not %#v", args))
			}
			return &Apply{r, m.Name, args[0]}
		}
	case *mast.Var:
		if unicode.IsDigit(rune(e.Name[0])) {
			v, _ := strconv.ParseInt(e.Name, 10, 64)
			return Const(v)
		} else if e.Name == "[]" {
			return EmptyList{}
		} else {
			for i, arg := range *args {
				if arg == e.Name {
					return &Var{i}
				}
			}
			if lval {
				*args = append(*args, e.Name)
				return &Var{len(*args) - 1}
			} else {
				return &Apply{r, e.Name, Const(0)}
			}
		}
	default:
		panic(fmt.Sprintf("not sure what to do with %v", e))
	}
}

func (r *Runtime) Parse(text string) error {
	if r.Funcs == nil {
		r.Funcs = map[string]Node{}
	}
	tree, err := parser.Parse(text)
	if err != nil {
		return err
	}

	names := []string{}
	args := []Node{}
	name := ""
	switch lhs := tree.Left.(type) {
	case *mast.Apply:
		name = lhs.Operator.(*mast.Var).Name
		args = r.mastToTuple(lhs.Operand, true, &names)
	case *mast.Var:
		name = lhs.Name
	default:
		return fmt.Errorf("not sure what to do with %s", lhs)
	}

	previous, ok := r.Funcs[name]
	if !ok {
		previous = &Undef{"failure to pattern match"}
	}

	rhs := r.mastToExpr(tree.Right, false, &names)
	for i, arg := range args {
		if _, ok := arg.(*Var); !ok {
			rhs = &If{&Plus{&Var{i}, &Negate{arg}}, rhs, previous}
		}
	}

	r.Funcs[name] = rhs
	return nil
}

func (r *Runtime) ParseFile(text string) error {
	lines := strings.Split(text, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if idx := strings.Index(line, "--"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if err := r.Parse(line); err != nil {
			return err
		}
	}
	return nil
}
