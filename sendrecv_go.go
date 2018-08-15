// +build ignore

package main

import "fmt"

func main() {
	left := func(ch chan struct{}) {
		fmt.Println("left: sending a to ch")
		a := struct{}{}
		ch <- a
	}
	right := func(ch chan struct{}) struct{} {
		x := <-ch
		fmt.Println("right: received x from ch")
		return x
	}

	ch := make(chan struct{}) // (new ch)
	go left(ch)               //   (  ch<a> |
	right(ch)                 //      ch(x).0 )
}
