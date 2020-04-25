package main

import (
	"flag"
	"fmt"
	"github.com/care0717/gogo/token"
	"log"

)

func printHeader() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global _main")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalln("引数の個数が正しくありません")
		return
	}
	t, err := token.Tokenize(args[0])
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	n := t.Expr()
	printHeader()
	fmt.Println("_main:")
	fmt.Println(n.Gen())
	fmt.Println("  pop rax")
	fmt.Println("  ret")
}
