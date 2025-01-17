package myerrors

import (
	"fmt"
	"os"
)

func Try(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	return
}

func TryWithValue[T any](value T, err error) T {
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	return value
}
