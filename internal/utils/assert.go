package utils

import (
	"fmt"
)

func Assert(condition bool, message string) {
	if !condition {
		panic(fmt.Sprintf("Assertion failed: %s", message))
	}
}
