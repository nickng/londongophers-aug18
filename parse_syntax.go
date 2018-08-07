// +build ignore

package main

import (
	"fmt"
	"log"
	"strings"

	"go.nickng.io/asyncpi"
)

func main() {
	const example = `
(new ch)    # Create a new channel ch
(
  ch<a>     # Send on channel ch with a as parameter
  |         # Parallel composition - think spawn goroutine
  ch(x)     # Receive x on channel ch
   .0       # End of process (receive always have continuation)
)`

	process, err := asyncpi.Parse(strings.NewReader(example))
	if err != nil {
		log.Fatal("parse failed:", err)
	}
	fmt.Println("AST:", process) // AST of process
	fmt.Println("Calculi:", process.Calculi())
}
