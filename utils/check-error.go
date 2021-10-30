package utils

import "log"

// Check - Function to log and exit the program accepting an error type.
func Check(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}
