package lox

import "fmt"

type AstPrinter struct{}

// AstPrinter implements visitorString
var _ visitorExprString = &AstPrinter{}

func (a *AstPrinter) Println(expr Expr) {
	fmt.Println(a.Sprint(expr))
}

func (a *AstPrinter) Sprint(expr Expr) string {
	return expr.AcceptString(a)
}

func (a *AstPrinter) visitBinaryExpr(expr *BinaryExpr) string {
	return a.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (a *AstPrinter) visitGroupingExpr(expr *GroupingExpr) string {
	return a.parenthesize("group", expr.expression)
}

func (a *AstPrinter) visitLiteralExpr(expr *LiteralExpr) string {
	if expr.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.value)
}

func (a *AstPrinter) visitUnaryExpr(expr *UnaryExpr) string {
	return a.parenthesize(expr.operator.Lexeme, expr.right)
}

func (a *AstPrinter) visitVariableExpr(expr *VariableExpr) string {
	return a.parenthesize("var", expr)
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
