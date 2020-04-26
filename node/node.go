package node

import (
	"fmt"
)

type kind int

const (
	Add kind = iota + 1
	Sub
	Mul
	Div
	Eq
	Ne
	Lt
	Le
	Num
)

type Node interface {
	Gen() string
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

func NewNumNode(val int) Node {
	return &node{
		kind: Num,
		val:  val,
	}
}

func (node *node) Gen() string {
	if node.kind == Num {
		return fmt.Sprintf("  push %d\n", node.val)
	}
	var res string

	res += node.lhs.Gen()
	res += node.rhs.Gen()
	res += fmt.Sprintln("  pop rdi")
	res += fmt.Sprintln("  pop rax")

	switch node.kind {
	case Add:
		res += fmt.Sprintln("  add rax, rdi")
	case Sub:
		res += fmt.Sprintln("  sub rax, rdi")
	case Mul:
		res += fmt.Sprintln("  imul rax, rdi")
	case Div:
		res += fmt.Sprintln("  cqo")
		res += fmt.Sprintln("  idiv rdi")
	case Eq:
		res += fmt.Sprintln("  cmp rax, rdi")
		res += fmt.Sprintln("  sete al")
		res += fmt.Sprintln("  movzx rax, al")
	case Ne:
		res += fmt.Sprintln("  cmp rax, rdi")
		res += fmt.Sprintln("  setne al")
		res += fmt.Sprintln("  movzx rax, al")
	case Lt:
		res += fmt.Sprintln("  cmp rax, rdi")
		res += fmt.Sprintln("  setl al")
		res += fmt.Sprintln("  movzx rax, al")
	case Le:
		res += fmt.Sprintln("  cmp rax, rdi")
		res += fmt.Sprintln("  setle al")
		res += fmt.Sprintln("  movzx rax, al")
	}
	res += fmt.Sprintln("  push rax")
	return res
}

