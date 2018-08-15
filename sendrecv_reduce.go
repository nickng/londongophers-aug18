// +build ignore

package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go.nickng.io/asyncpi"
	"golang.org/x/tools/imports"
)

func main() {
	const sendrecv = `
	(new ch)(       # Create a channel "ch"
	  ch<a>         # Send to "ch"
	  | ch(x).0     # Concurrently, Receive from "ch"
	)`
	proc, err := asyncpi.Parse(strings.NewReader(sendrecv))
	if err != nil {
		// handle error
	}
	fmt.Printf("Before reduction:\n\t%s\n", proc.Calculi())

	// Reduce system for a single step // HL
	asyncpi.Reduce1(proc) // HL

	fmt.Printf("After reduction:\n\t%s\n", proc.Calculi())
}

func mustParse(s string) asyncpi.Process {
	proc, err := asyncpi.Parse(strings.NewReader(s))
	if err != nil {
		log.Fatal("parse failed:", err)
	}
	return proc
}

func reduceAll(proc asyncpi.Process) asyncpi.Process {
	for {
		changed, err := asyncpi.Reduce1(proc)
		if err != nil {
			log.Fatal("reduction error", err) // handle errors
			break
		}
		if !changed {
			break
		}
		proc, _ = asyncpi.SimplifyBySC(proc)
		fmt.Println("→ Reduces to:", proc.Calculi())
	}
	return proc
}

func fixImports(src []byte) []byte {
	opts := &imports.Options{TabIndent: true, Fragment: false}
	imported, err := imports.Process("/tmp/main.go", src, opts)
	if err != nil {
		log.Fatal(err)
	}
	return imported
}

func goFmt(src []byte) []byte {
	fmtd, err := format.Source(src)
	if err != nil {
		log.Fatal(err)
	}
	return fmtd
}

func writeTemp(content []byte) {
	f, err := ioutil.TempFile("", "generated")
	os.Rename(f.Name(), f.Name()+".go")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		log.Fatal(err)
	}
	fmt.Println("written to temp file:", f.Name()+".go")
}
