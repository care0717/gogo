package token

import (
	"fmt"
	"github.com/care0717/gogo/node"
	"regexp"
	"strconv"
)

type Token interface {
	Program() ([]node.Node, error)
}

func (t *token) Program() ([]node.Node, error) {
	var code []node.Node
	for !t.atEof() {
		c, err := t.stmt()
		if err != nil {
			return nil, err
		}
		code = append(code, c)
	}
	return code, nil
}

func (t *token) stmt() (node.Node, error) {
	n, err := t.expr()
	if err != nil {
		return nil, err
	}
	if err := t.expect(";"); err != nil {
		return nil, err
	}
	return n, nil
}

func (t *token) expr()  (node.Node, error) {
	return t.assign()
}

func (t *token) assign()  (node.Node, error) {
	n, err := t.equality()
	if err != nil {
		return nil, err
	}
	if t.consume("=") {
		tmp, err := t.assign()
		if err != nil {
			return nil, err
		}
		n = node.NewNode(node.Assign, n, tmp)
	}
	return n, nil
}

func (t *token) equality() (node.Node, error) {
	n, err:= t.relational()
	if err != nil {
		return nil, err
	}
	for {
		if t.consume("==") {
			tmp, err:= t.relational()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Eq, n, tmp)
		} else if t.consume("!=") {
			tmp, err:= t.relational()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Ne, n, tmp)
		} else {
			return n, nil
		}
	}
}

func (t *token) relational() (node.Node, error) {
	n, err := t.add()
	if err != nil {
		return nil, err
	}
	for {
		if t.consume("<") {
			tmp, err := t.add()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Lt, n, tmp)
		} else if t.consume(">") {
			tmp, err := t.add()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Lt, tmp, n)
		} else if t.consume("<=") {
			tmp, err := t.add()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Le, n, tmp)
		} else if t.consume(">=") {
			tmp, err := t.add()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Le, tmp, n)
		} else {
			return n, nil
		}
	}
}

func (t *token) add() (node.Node, error) {
	n, err := t.mul()
	if err != nil {
		return nil, err
	}
	for {
		if t.consume("+") {
			tmp, err := t.mul()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Add, n, tmp)
		} else if t.consume("-") {
			tmp, err := t.mul()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Sub, n, tmp)
		} else {
			return n, nil
		}
	}
}

func (t *token) mul() (node.Node, error) {
	n, err := t.unary()
	if err != nil {
		return nil, err
	}
	for {
		if t.consume("*") {
			tmp, err := t.unary()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Mul, n, tmp)
		} else if t.consume("/") {
			tmp, err := t.unary()
			if err != nil {
				return nil, err
			}
			n = node.NewNode(node.Div, n, tmp)
		} else {
			return n, nil
		}
	}
}

func (t *token) unary() (node.Node, error) {
	if t.consume("+") {
		return t.unary()
	} else if t.consume("-") {
		n, err := t.unary()
		if err != nil {
			return nil, err
		}
		return node.NewNode(node.Sub, node.NewNumNode(0), n), nil
	}
	return t.primary()
}

func (t *token) primary() (node.Node, error) {
	if t.consume("(") {
		n, err := t.expr()
		if err != nil {
			return nil, err
		}
		if err := t.expect(")"); err != nil {
			return nil, err
		}
		return n, nil
	}

	if i, ok := t.consumeNumber(); ok {
		return node.NewNumNode(i), nil
	}
	if s, ok := t.consumeIdent(); ok {
		if lvar, ok := lVarMap[s]; ok {
			return node.NewLVarNode(lvar.offset), nil
		}
		offset += 8
		lVar := LVar{
			name:   s,
			offset: offset,
		}
		lVarMap[s] = lVar
		return node.NewLVarNode(offset), nil
	}
	return nil, compileError{fmt.Sprintf("不明な識別子 %s", t.str), t.meta.line, t.meta.pos}
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
	ident
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
		return compileError{fmt.Sprintf("%sではありません", op), t.meta.line, t.meta.pos}
	}
	*t = *t.next
	return nil
}

func (t *token) consumeNumber() (int, bool) {
	if t.kind != number {
		return 0, false
	}
	val := t.value
	*t = *t.next
	return val, true
}

func (t *token) consumeIdent() (string, bool) {
	if t.kind != ident {
		return "", false
	}
	val := t.str
	*t = *t.next
	return val, true
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
var regexAlphabet = regexp.MustCompile(`^[a-zA-Z]+`)
var regexOp = regexp.MustCompile(`^(\+|-|\*|/|\(|\)|=|;|<=|<|>|>=|==|!=)`)

type LVar struct {
	name string
	offset int
}

var lVarMap = map[string]LVar{}
var offset int

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
		if regexAlphabet.MatchString(s[i:]) {
			variable := regexAlphabet.FindString(s[i:])
			cur = newNextToken(ident, cur, variable, meta{s, i})
			i += len(variable)
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
