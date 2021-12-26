package glox

import "fmt"

type astPrinter struct{}

// astPrinter implements visitorString
var _ visitorString = &astPrinter{}

func (a *astPrinter) Print(expr Expr) string {
	return expr.AcceptString(a)
}

func (a *astPrinter) visitBinaryExpr(expr *Binary) string {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (a *astPrinter) visitGroupingExpr(expr *Grouping) string {
	return a.parenthesize("group", expr.expression)
}

func (a *astPrinter) visitLiteralExpr(expr *Literal) string {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (a *astPrinter) visitUnaryExpr(expr *Unary) string {
	return a.parenthesize(expr.operator.Lexeme, expr.right)
}

func (a *astPrinter) parenthesize(name string, exprs ...Expr) string {
	s := fmt.Sprintf("(%s", name)
	for _, expr := range exprs {
		s += " "
		s += a.Print(expr)
	}
	s += ")"
	return s
}
