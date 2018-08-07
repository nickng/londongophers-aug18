//go:generate goyacc -p asyncpi -o parser.y.go asyncpi.y

package main

func main() {
	// ... uses parser.y.go
}
