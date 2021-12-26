package glox

import "fmt"

type AstPrinter struct{}

// astPrinter implements visitorString
var _ visitorString = &AstPrinter{}

func (a *AstPrinter) Println(expr Expr) {
	fmt.Println(a.Sprint(expr))
}

func (a *AstPrinter) Sprint(expr Expr) string {
	return expr.AcceptString(a)
}

func (a *AstPrinter) visitBinaryExpr(expr *Binary) string {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitGroupingExpr(expr *Grouping) string {
	return a.parenthesize("group", expr.expression)
}

func (a *AstPrinter) visitLiteralExpr(expr *Literal) string {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (a *AstPrinter) visitUnaryExpr(expr *Unary) string {
	return a.parenthesize(expr.operator.Lexeme, expr.right)
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	s := fmt.Sprintf("(%s", name)
	for _, expr := range exprs {
		s += " "
		s += a.Sprint(expr)
	}
	s += ")"
	return s
}
