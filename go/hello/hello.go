package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {
	// Set properties of the predifiend logger, including
	// the log entry prefix and a flag to disable printing
	// the time, source file, and line number
	log.SetPrefix("greetings: ")
	log.SetFlags(0)

	names := []string{"Kait", "Luciana", "Ezra", "Pato"}
	// Request a greeting message
	message, err := greetings.Hellos(names)
	// if an error was returned, print it to the console and
	// exit the program
	if err != nil {
		log.Fatal(err)
	}
	// if no error was returned, print the returned map of
	// messages to the console.
	fmt.Println(message)
}
