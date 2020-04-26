package main

import (
	"bufio"
	"fmt"
	"github.com/care0717/gogo/token"
	"log"
	"os"
)

func printHeader() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global _main")
}

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	stdin.Scan()
	text := stdin.Text()
	t, err := token.Tokenize(text)
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
