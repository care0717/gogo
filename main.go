package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type TokenKind int

type CompileError struct {
	message string
	line    string
	pos     int
}

type Meta struct {
	line string
	pos  int
}

func (c CompileError) Error() string {
	var err string
	err += fmt.Sprintf("\n%s\n", c.line)
	err += fmt.Sprintf("%*s", c.pos, " ")
	err += fmt.Sprintf("^ ")
	err += fmt.Sprintf("%s", c.message)
	return err
}

const (
	TK_RESERVED TokenKind = iota
	TK_NUM
	TK_EOF
)

type Token struct {
	kind  TokenKind
	next  *Token
	value int
	str   string
	meta  Meta
}

func (t *Token) consume(op byte) bool {
	if t.kind != TK_RESERVED || t.str[0] != op {
		return false
	}
	*t = *t.next
	return true
}
func (t *Token) expect(op byte) error {
	if t.kind != TK_RESERVED || t.str[0] != op {
		return CompileError{fmt.Sprintf("%cではありません", op), t.meta.line, t.meta.pos}
	}
	*t = *t.next
	return nil
}

func (t *Token) expectNumber() (int, error) {
	if t.kind != TK_NUM {
		return 0, CompileError{"数ではありません", t.meta.line, t.meta.pos}
	}
	val := t.value
	*t = *t.next
	return val, nil
}

func (t Token) atEof() bool {
	return t.kind == TK_EOF
}

func newNextToken(kind TokenKind, cur *Token, str string, meta Meta) *Token {
	tok := &Token{
		kind: kind,
		str:  str,
		meta: meta,
	}
	cur.next = tok
	return tok
}

var regexNumber = regexp.MustCompile(`^[0-9]+`)

func tokenize(s string) (*Token, error) {
	head := Token{}
	cur := &head
	i := 0
	for i < len(s) {
		if s[i] == ' ' {
			i++
			continue
		}
		if s[i] == '+' || s[i] == '-' {
			cur = newNextToken(TK_RESERVED, cur, string(s[i]), Meta{s, i})
			i++
			continue
		}
		if regexNumber.MatchString(s[i:]) {
			num := regexNumber.FindString(s[i:])
			cur = newNextToken(TK_NUM, cur, num, Meta{s, i})
			n, _ := strconv.Atoi(num)
			cur.value = n
			i += len(num)
			continue
		}

		return nil, CompileError{"tokenizeできません", s, i}
	}
	newNextToken(TK_EOF, cur, "", Meta{s, i})
	return head.next, nil
}

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
	token, err := tokenize(args[0])
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	printHeader()
	fmt.Println("_main:")
	n, err := token.expectNumber()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Printf("  mov rax, %d\n", n)
	for !token.atEof() {
		if token.consume('+') {
			n, err := token.expectNumber()
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			fmt.Printf("  add rax, %d\n", n)
			continue
		}
		if token.consume('-') {
			n, err := token.expectNumber()
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			fmt.Printf("  sub rax, %d\n", n)
			continue
		}
	}
	fmt.Println("  ret")
}
