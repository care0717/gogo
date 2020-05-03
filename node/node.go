package node

import (
	"errors"
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
	Assign
	Lvar
	Num
)

type Node interface {
	Gen() string
	genLval() (string, error)
}

type node struct {
	kind kind
	lhs  Node
	rhs  Node
	val  int
	offset int
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

func NewLVarNode(offset int) Node {
	return &node{
		kind: Lvar,
		offset:  offset,
	}
}

func (node *node) Gen() string {
	switch node.kind  {
	case Num:
		return fmt.Sprintf("  push %d\n", node.val)
	case Lvar:
		res, _ := node.genLval()
		res += fmt.Sprintln("  pop rax")
		res += fmt.Sprintln("  mov rax, [rax]")
		res += fmt.Sprintln("  push rax")
		return res
	case Assign:
		res, _ := node.lhs.genLval()
		res += node.rhs.Gen()
		res += fmt.Sprintln("  pop rdi")
		res += fmt.Sprintln("  pop rax")
		res += fmt.Sprintln("  mov [rax], rdi")
		res += fmt.Sprintln("  push rdi")
		return res
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

func (node *node) genLval() (string, error) {
	if node.kind != Lvar {
		return "", errors.New("左辺が変数ではありません")
	}
	var res string
	res += fmt.Sprintln("  mov rax, rbp")
	res += fmt.Sprintf("  sub rax, %d\n", node.offset)
	res += fmt.Sprintln("  push rax")
	return res, nil
}

