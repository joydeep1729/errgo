package main

import (
	"fmt"

	errors "github.com/joydeep1729/errgo"
)

func main() {
	// 1. Create a nested error chain:
	// dbErr (root) -> appErr (middle context) -> apiErr (top context)
	dbErr := errors.New("database connection failed")
	appErr := errors.WithMessage(dbErr, "failed to fetch user profile")
	apiErr := errors.WithMessage(appErr, "API request failed")

	fmt.Println("=== Original Wrapped Error ===")
	fmt.Printf("%v\n\n", apiErr)

	// 2. Demonstrate UnwrapWithOuter (single-level unwrap)
	fmt.Println("=== Step-by-Step unwrapping using UnwrapWithOuter() ===")
	current := apiErr
	step := 1
	for {
		inner, outer := errors.UnwrapWithOuter(current)
		if inner == nil {
			fmt.Printf("Step %d (Leaf): Outer context = %v | Inner = <nil>\n", step, outer)
			break
		}
		fmt.Printf("Step %d: Outer context = %q | Inner error msg = %q\n", step, outer.Error(), inner.Error())
		current = inner
		step++
	}
	fmt.Println()

	// 3. Demonstrate UnwrapToCauseWithOuter (multi-level unwrap to root cause)
	fmt.Println("=== Direct unwrapping to root cause using UnwrapToCauseWithOuter() ===")
	cause, outer := errors.UnwrapToCauseWithOuter(apiErr)
	fmt.Printf("Root Cause  : %v\n", cause)
	fmt.Printf("Outer Context: %v\n", outer)
}
