package token

import (
	"fmt"
	"github.com/care0717/gogo/node"
	"regexp"
	"strconv"
)

type Token interface {
	Expr() node.Node
}

func (t *token) Expr() node.Node {
	return t.equality()
}

func (t *token) equality() node.Node {
	n := t.relational()
	for {
		if t.consume("==") {
			n = node.NewNode(node.Eq, n, t.relational())
		} else if t.consume("!=") {
			n = node.NewNode(node.Ne, n, t.relational())
		} else {
			return n
		}
	}
}

func (t *token) relational() node.Node {
	n := t.add()
	for {
		if t.consume("<") {
			n = node.NewNode(node.Lt, n, t.add())
		} else if t.consume(">") {
			n = node.NewNode(node.Lt, t.add(), n)
		} else if t.consume("<=") {
			n = node.NewNode(node.Le, n, t.add())
		} else if t.consume(">=") {
			n = node.NewNode(node.Le, t.add(), n)
		} else {
			return n
		}
	}
}

func (t *token) add() node.Node {
	n := t.mul()
	for {
		if t.consume("+") {
			n = node.NewNode(node.Add, n, t.mul())
		} else if t.consume("-") {
			n = node.NewNode(node.Sub, n, t.mul())
		} else {
			return n
		}
	}
}

func (t *token) mul() node.Node {
	n := t.unary()
	for {
		if t.consume("*") {
			n = node.NewNode(node.Mul, n, t.unary())
		} else if t.consume("/") {
			n = node.NewNode(node.Div, n, t.unary())
		} else {
			return n
		}
	}
}

func (t *token) unary() node.Node {
	if t.consume("+") {
		return t.unary()
	} else if t.consume("-") {
		return node.NewNode(node.Sub, node.NewNumNode(0), t.unary())
	}
	return t.primary()
}

func (t *token) primary() node.Node {
	if t.consume("(") {
		n := t.Expr()
		t.expect(")")
		return n
	}
	i, _ := t.expectNumber()
	return node.NewNumNode(i)
}

type compileError struct {
	message string
	line    string
	pos     int
}

type meta struct {
	line string
	pos  int
}

func (c compileError) Error() string {
	var err string
	err += fmt.Sprintf("\n%s\n", c.line)
	err += fmt.Sprintf("%*s", c.pos, " ")
	err += fmt.Sprintf("^ ")
	err += fmt.Sprintf("%s", c.message)
	return err
}

type kind int

const (
	reserved kind = iota + 1
	number
	eof
)

type token struct {
	kind  kind
	next  *token
	value int
	str   string
	meta  meta
}

func (t *token) consume(op string) bool {
	if t.kind != reserved || t.str != op {
		return false
	}
	*t = *t.next
	return true
}
func (t *token) expect(op string) error {
	if t.kind != reserved || t.str != op {
		return compileError{fmt.Sprintf("%cではありません", op), t.meta.line, t.meta.pos}
	}
	*t = *t.next
	return nil
}

func (t *token) expectNumber() (int, error) {
	if t.kind != number {
		return 0, compileError{"数ではありません", t.meta.line, t.meta.pos}
	}
	val := t.value
	*t = *t.next
	return val, nil
}

func (t token) atEof() bool {
	return t.kind == eof
}

func newNextToken(kind kind, cur *token, str string, meta meta) *token {
	tok := &token{
		kind: kind,
		str:  str,
		meta: meta,
	}
	cur.next = tok
	return tok
}

var regexNumber = regexp.MustCompile(`^[0-9]+`)
var regexOp = regexp.MustCompile(`^(\+|-|\*|/|\(|\)|<=|<|>|>=|==|!=)`)

func Tokenize(s string) (Token, error) {
	head := token{}
	cur := &head
	i := 0
	regexOp.Longest()
	for i < len(s) {
		if s[i] == ' ' {
			i++
			continue
		}
		if regexOp.MatchString(s[i:]) {
			op := regexOp.FindString(s[i:])
			cur = newNextToken(reserved, cur, op, meta{s, i})
			i += len(op)
			continue
		}
		if regexNumber.MatchString(s[i:]) {
			num := regexNumber.FindString(s[i:])
			cur = newNextToken(number, cur, num, meta{s, i})
			n, err := strconv.Atoi(num)
			if err != nil {
				return nil, err
			}
			cur.value = n
			i += len(num)
			continue
		}

		return nil, compileError{"tokenizeできません", s, i}
	}
	newNextToken(eof, cur, "", meta{s, i})
	return head.next, nil
}
