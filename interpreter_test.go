package madison_test

import (
	"fmt"
	. "github.com/fatlotus/madison"
)

func ExampleNode_Eval() {
	r := &Runtime{}
	if err := r.ParseFile(`
		fib 0 = 1
		fib 1 = 1
		fib n = fib(n - 1) + fib(n - 2)
		
		repeat 0 = []
		repeat n = n : repeat(n - 1)

		safe = head(repeat 10)
		unsafe = head(repeat 0) 
	`); err != nil {
		panic(err)
	}

	// Compute how large the sixth Fibbonacci number is.
	fib5 := r.Funcs["fib"].Eval([]Obj{Obj{Int: 5}})
	fmt.Printf("fib 5 = %s\n", fib5)

	// Create a list with a fixed range of values
	repeat := r.Funcs["repeat"].Eval([]Obj{Obj{Int: 3}})
	fmt.Printf("repeat 3 = %s\n", repeat)

	// Output:
	// fib 5 = 8
	// repeat 3 = 3 : 2 : 1 : []
}
