package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalln("引数の個数が正しくありません")
		return
	}
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".global _main")
	fmt.Println("_main:")
	fmt.Printf("  mov rax, %s\n", args[0])
	fmt.Println("  ret")
}
