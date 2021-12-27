package lox

import "fmt"

type parser struct {
	tokens  []Token
	current int

	errors []*parseError
}

// NewParser creates a parser for the lox language. The complete expression grammar is the following:
//
// program    → statement* EOF ;
//
// statement  → exprStmt | printStmt ;
//
// exprStmt   → expression ";" ;
//
// printStmt  → "print" expression ";" ;
//
// expression → equality ;
//
// equality   → comparison ( ( "!=" | "==" ) comparison )* ;
//
// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
//
// term       → factor ( ( "-" | "+" ) factor )* ;
//
// factor     → unary ( ( "/" | "*" ) unary )* ;
//
// unary      → ( "!" | "-" ) unary | primary ;
//
// primary    → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func NewParser(tokens []Token) *parser {
	return &parser{
		tokens:  tokens,
		current: 0,
		errors:  make([]*parseError, 0),
	}
}

func (p *parser) Parse() Expr {
	expr, err := p.expression()
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	return expr
}

// TODO change this to see how we want to handle errors
func (p *parser) HadErrors() bool {
	hadErrors := false
	for _, err := range p.errors {
		println(err.Error())
		hadErrors = true
	}
	return hadErrors
}

func (p *parser) expression() (Expr, *parseError) {
	return p.equality()
}

func (p *parser) equality() (Expr, *parseError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) comparison() (Expr, *parseError) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) term() (Expr, *parseError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(Plus, Minus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) factor() (Expr, *parseError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(Star, Slash) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) unary() (Expr, *parseError) {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnaryExpr(operator, right), nil
	}
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}
	return expr, nil
}

func (p *parser) primary() (Expr, *parseError) {
	if p.match(False) {
		return NewLiteralExpr(false), nil
	}
	if p.match(True) {
		return NewLiteralExpr(true), nil
	}
	if p.match(Nil) {
		return NewLiteralExpr(nil), nil
	}

	if p.match(Number, String) {
		return NewLiteralExpr(p.previous().Literal), nil
	}

	if p.match(LeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewGroupingExpr(expr), nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *parser) consume(t TokenType, message string) (Token, *parseError) {
	if p.check(t) {
		return p.advance(), nil
	}
	return Token{}, p.error(p.peek(), message)
}

func (p *parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *parser) peek() Token {
	return p.tokens[p.current]
}

func (p *parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *parser) error(token Token, message string) *parseError {
	where := "at end"
	if token.Type != EOF {
		where = fmt.Sprintf("at '%s'", token.Lexeme)
	}
	return &parseError{line: token.Line, where: where, message: message}
}

func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}
		switch p.peek().Type {
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}
		p.advance()
	}
}

type parseError struct {
	line    int
	where   string
	message string
}

func (e *parseError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %s", e.line, e.where, e.message)
}
