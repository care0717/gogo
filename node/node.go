package node

import (
	"fmt"
	"log"
)

type kind int

const (
	Add kind = iota + 1
	Sub
	Mul
	Div
	Num
)

type Node interface {
	Gen() string
	debug()
}

type node struct {
	kind kind
	lhs  Node
	rhs  Node
	val  int
}

func NewNode(kind kind, lhs Node, rhs Node) Node {
	return &node{
		kind: kind,
		lhs:  lhs,
		rhs:  rhs,
	}
}

func NewNumNode(val int) Node{
	return &node{
		kind: Num,
		val:  val,
	}
}


func (node *node) Gen() string {
	if node.kind == Num {
		return fmt.Sprintf("  push %d\n", node.val)
	}
	var result string

	result += node.lhs.Gen()
	result += node.rhs.Gen()
	result += fmt.Sprintln("  pop rdi")
	result += fmt.Sprintln("  pop rax")

	switch node.kind {
	case Add:
		result +=fmt.Sprintln("  add rax, rdi")
	case Sub:
		result +=fmt.Sprintln("  sub rax, rdi")
	case Mul:
		result +=fmt.Sprintln("  imul rax, rdi")
	case Div:
		result +=fmt.Sprintln("  cqo")
		result +=fmt.Sprintln("  idiv rdi")
	}
	result += fmt.Sprintln("  push rax")
	return result
}

func (node *node) debug() {
	log.Println(node)
	if node.lhs != nil {
		node.lhs.debug()
	}
	if node.rhs != nil {
		node.rhs.debug()
	}
}
