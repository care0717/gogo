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
	fmt.Println("_main:")
	fmt.Println("  push rbp")
	fmt.Println("  mov rbp, rsp")
	fmt.Println("  sub rsp, 208")
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
	printHeader()
	ps, err := t.Program()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	for _, p := range ps {
		fmt.Println(p.Gen())
		fmt.Println("  pop rax")
	}
	fmt.Println("  mov rsp, rbp")
	fmt.Println("  pop rbp")
	fmt.Println("  ret")
}
