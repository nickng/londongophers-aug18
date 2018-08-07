// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go.nickng.io/asyncpi"
	"go.nickng.io/asyncpi/codegen/golang"
	"golang.org/x/tools/imports"
)

const chainReaction = `(new ch2)(ch2<1, 2> | ch2(x, y).((new ch3) (ch3<x> | ch3(z).ch3(z).0 | ch3<y>)))`
const sendrecv = `(new ch)(ch<> | ch().ch().0)`

var gocode bytes.Buffer

var withMain = golang.FormatOptions{
	Debug: true,
	Main:  true,
}

func main() {
	const srr = `(new ch)(ch<> | ch().ch().0)`
	proc := mustParse(srr)
	golang.GenerateOpts(proc, withMain, &gocode)
	fmt.Println(proc.Calculi())
	reduceAll(proc)

	generated := fixImports(gocode.Bytes())
	fmt.Println(string(generated))
	writeTemp(generated)
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
