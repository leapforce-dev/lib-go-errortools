package errortools

import (
	"fmt"
	"log"
)

// Println prints error if not nil
//
func Println(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// Fatal prints error and exits if not nil
//
func Fatal(err error) {
	if err != nil {
		fmt.Println("hallo")
		log.Fatal(err)
	}
}
