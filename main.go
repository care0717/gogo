package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

type NodeKind int

const (
	ND_ADD NodeKind = iota + 1
	ND_SUB
	ND_MUL
	ND_DIV
	ND_NUM
)

type Node struct {
	kind NodeKind
	lhs *Node
	rhs *Node
	val int
}

func NewNode(kind NodeKind, lhs *Node, rhs *Node) *Node {
	return &Node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func NewNumNode(val int) *Node{
	return &Node{
		kind: ND_NUM,
		val:  val,
	}
}

func (t *Token) expr() *Node {
	node := t.mul()
	for {
		if t.consume('+') {
			node = NewNode(ND_ADD, node, t.mul())
		} else if t.consume('-') {
			node = NewNode(ND_SUB, node, t.mul())
		} else {
			return node
		}
	}
}

func (t *Token) mul() *Node {
	node := t.primary()
	for {
		if t.consume('*') {
			node = NewNode(ND_MUL, node, t.primary())
		} else if t.consume('/') {
			node = NewNode(ND_DIV, node, t.primary())
		} else {
			return node
		}
	}
}

func (t *Token) primary() *Node {
	if t.consume('(') {
		node := t.expr()
		t.expect(')')
		return node
	}
	i, _ := t.expectNumber()
	return NewNumNode(i)
}

func (node *Node) debug() {
	log.Println(node)
	if node.lhs != nil {
		node.lhs.debug()
	}
	if node.rhs != nil {
		node.rhs.debug()
	}
}

func (node *Node) gen() {
	if node.kind == ND_NUM  {
		fmt.Printf("  push %d\n", node.val)
		return
	}

	node.lhs.gen()
	node.rhs.gen()
	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	switch node.kind {
	case ND_ADD:
		fmt.Println("  add rax, rdi")
	case ND_SUB:
		fmt.Println("  sub rax, rdi")
	case ND_MUL:
		fmt.Println("  imul rax, rdi")
	case ND_DIV:
		fmt.Println("  cqo")
		fmt.Println("  idiv rdi")
	}
	fmt.Println("  push rax")
}

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

type TokenKind int
const (
	TK_RESERVED TokenKind = iota + 1
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
var tokenRegexp = regexp.MustCompile(`([+\-*/()])`)

func tokenize(s string) (*Token, error) {
	head := Token{}
	cur := &head
	i := 0
	for i < len(s) {
		if s[i] == ' ' {
			i++
			continue
		}
		if tokenRegexp.Match([]byte{s[i]}) {
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
	node := token.expr()
	printHeader()
	fmt.Println("_main:")
	node.gen()
	fmt.Println("  pop rax")
	fmt.Println("  ret")
}
