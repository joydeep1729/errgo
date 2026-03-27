package main

import (
	"fmt"

	errors "errgo"
)

func main() {
	// Create a base error simulating a database failure
	dbErr := errors.New("database connection failed")

	// Wrap the error with custom application-level context
	appErr := errors.WithMessage(dbErr, "failed to fetch user profile")

	// Demonstrate our new UnwrapWithOuter feature
	inner, outer := errors.UnwrapWithOuter(appErr)

	fmt.Println("--- Original Wrapped Error ---")
	fmt.Printf("%v\n\n", appErr)

	fmt.Println("--- After UnwrapWithOuter() ---")
	fmt.Printf("Inner Error (Original Cause): %v\n", inner)
	fmt.Printf("Outer Error (The Wrapper)   : %v\n", outer)
}
