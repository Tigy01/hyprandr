package myerrors

import (
	"fmt"
	"os"
)

// Calls Print and os.Exit(1)
func Try(err error) {
	if err != nil {
        fmt.Printf("Err: %v\n", err)
		os.Exit(1)
	}
	return
}

// Returns the value if err == nil
func TryWithValue[T any](value T, err error) T {
	if err != nil {
        fmt.Printf("Err: %v\n", err)
		os.Exit(1)
	}
	return value
}
