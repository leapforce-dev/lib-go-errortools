package errortools

import (
	"fmt"
	"log"
)

// Println prints error if not nil
//
func Println(a ...interface{}) {
	if a != nil {
		fmt.Println(a...)
	}
}

// Fatal prints error and exits if not nil
//
func Fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
