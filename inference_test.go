package madison_test

import (
	"fmt"
	. "github.com/fatlotus/madison"
)

func Example() {
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

	// Compute how large the first 5 Fibbonacci numbers are.
	typ, _ := r.Funcs["fib"].Type(nil, []Type{InRange(0, 5)})
	fmt.Printf("fib :: [0, 5] -> %s\n", typ)

	// Create a list with a fixed range of values
	typ, _ = r.Funcs["repeat"].Type(nil, []Type{InRange(3, 5)})
	fmt.Printf("repeat :: [3, 5] -> %s\n", typ)

	// Make sure there aren't unsafe head() calls.
	_, err := r.Funcs["unsafe"].Type(nil, []Type{})
	fmt.Printf("unsafe raises %s\n", err)
	// Output:
	// fib :: [0, 5] -> int[1, 8]
	// repeat :: [3, 5] -> [3, 5]int[1, 5]
	// unsafe raises cannot take head of an empty list: [0]any
}
