package main

import (
	"fmt"
	"go/format"
	"log"
	"strings"

	"go.nickng.io/asyncpi"
	"golang.org/x/tools/imports"
)

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
