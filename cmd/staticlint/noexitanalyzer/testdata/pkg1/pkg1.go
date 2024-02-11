package main

import (
	"fmt"
	"os"
)

func helloFunc() {
	fmt.Printf("Hello world")
}

func main() {
	helloFunc()
	os.Exit(0) // want "using os.Exit in main func is forbidden"
	helloFunc()
}
