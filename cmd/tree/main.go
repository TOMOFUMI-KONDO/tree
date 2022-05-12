package main

import (
	"flag"
	"log"

	"github.com/TOMOFUMI-KONDO/tree"
)

var path string

func init() {
	flag.StringVar(&path, "path", "./", "path of tree root")
	flag.Parse()
}

func main() {
	if err := tree.PrintTree(path); err != nil {
		log.Fatalln(err)
	}
}
