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
	n := t.mul()
	for {
		if t.consume('+') {
			n = node.NewNode(node.ND_ADD, n, t.mul())
		} else if t.consume('-') {
			n = node.NewNode(node.ND_SUB, n, t.mul())
		} else {
			return n
		}
	}
}

func (t *token) mul() node.Node {
	n := t.primary()
	for {
		if t.consume('*') {
			n = node.NewNode(node.ND_MUL, n, t.primary())
		} else if t.consume('/') {
			n = node.NewNode(node.ND_DIV, n, t.primary())
		} else {
			return n
		}
	}
}

func (t *token) primary() node.Node {
	if t.consume('(') {
		n := t.Expr()
		t.expect(')')
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
	TK_RESERVED kind = iota + 1
	TK_NUM
	TK_EOF
)

type token struct {
	kind  kind
	next  *token
	value int
	str   string
	meta  meta
}

func (t *token) consume(op byte) bool {
	if t.kind != TK_RESERVED || t.str[0] != op {
		return false
	}
	*t = *t.next
	return true
}
func (t *token) expect(op byte) error {
	if t.kind != TK_RESERVED || t.str[0] != op {
		return compileError{fmt.Sprintf("%cではありません", op), t.meta.line, t.meta.pos}
	}
	*t = *t.next
	return nil
}

func (t *token) expectNumber() (int, error) {
	if t.kind != TK_NUM {
		return 0, compileError{"数ではありません", t.meta.line, t.meta.pos}
	}
	val := t.value
	*t = *t.next
	return val, nil
}

func (t token) atEof() bool {
	return t.kind == TK_EOF
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
var tokenRegexp = regexp.MustCompile(`([+\-*/()])`)

func Tokenize(s string) (Token, error) {
	head := token{}
	cur := &head
	i := 0
	for i < len(s) {
		if s[i] == ' ' {
			i++
			continue
		}
		if tokenRegexp.Match([]byte{s[i]}) {
			cur = newNextToken(TK_RESERVED, cur, string(s[i]), meta{s, i})
			i++
			continue
		}
		if regexNumber.MatchString(s[i:]) {
			num := regexNumber.FindString(s[i:])
			cur = newNextToken(TK_NUM, cur, num, meta{s, i})
			n, _ := strconv.Atoi(num)
			cur.value = n
			i += len(num)
			continue
		}

		return nil, compileError{"tokenizeできません", s, i}
	}
	newNextToken(TK_EOF, cur, "", meta{s, i})
	return head.next, nil
}
