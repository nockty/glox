// Code generated by the generate_ast tool; DO NOT EDIT.
package glox

type Expr interface {
	accept(v Visitor) interface{}
}

type Visitor interface {
	visitBinary(*Binary) interface{}
	visitGrouping(*Grouping) interface{}
	visitLiteral(*Literal) interface{}
	visitUnary(*Unary) interface{}
}

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func (expr *Binary) accept(v Visitor) interface{} {
	return v.visitBinary(expr)
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

type Grouping struct {
	expression Expr
}

func (expr *Grouping) accept(v Visitor) interface{} {
	return v.visitGrouping(expr)
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{
		expression: expression,
	}
}

type Literal struct {
	value interface{}
}

func (expr *Literal) accept(v Visitor) interface{} {
	return v.visitLiteral(expr)
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{
		value: value,
	}
}

type Unary struct {
	operator Token
	right    Expr
}

func (expr *Unary) accept(v Visitor) interface{} {
	return v.visitUnary(expr)
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{
		operator: operator,
		right:    right,
	}
}
