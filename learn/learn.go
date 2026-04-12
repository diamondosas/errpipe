package main

import (
	"flag"
	"fmt"
	"strings"
)

func greetUser(name string, upper bool) {
	msg := "Hello, " + name + "!"
	if upper {
		msg = strings.ToUpper(msg)
	}
	fmt.Println(msg)
}

func addNumbers(a, b int) {
	fmt.Printf("%d + %d = %d\n", a, b, a+b)
}

func repeatText(text string, times int) {
	for i := 0; i < times; i++ {
		fmt.Println(text)
	}
}

func main() {
	// --- Define flags ---
	// String flag
	name := flag.String("name", "World", "Name to greet")

	// Bool flag
	upper := flag.Bool("upper", false, "Uppercase the greeting")

	// Int flags
	numA := flag.Int("a", 0, "First number")
	numB := flag.Int("b", 0, "Second number")

	// Repeat flags																								
	text  := flag.String("text", "", "Text to repeat")
	times := flag.Int("times", 1, "How many times to repeat")

	// Mode flag — used to pick which function to run
	mode := flag.String("mode", "greet", "Mode: greet | add | repeat")

	flag.Parse()

	// --- Trigger different functions based on -mode ---
	switch *mode {
	case "greet":
		greetUser(*name, *upper)

	case "add":
		addNumbers(*numA, *numB)

	case "repeat":
		if *text == "" {
			fmt.Println("Error: -text is required in repeat mode")
			return
		}
		repeatText(*text, *times)

	default:
		fmt.Printf("Unknown mode: %q\n", *mode)
		fmt.Println("Available modes: greet | add | repeat")
	}
}