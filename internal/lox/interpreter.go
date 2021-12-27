package lox

import "fmt"

type Interpreter struct{}

// Interpreter implements visitorExpr and visitorStmt
var _ visitorExpr = &Interpreter{}
var _ visitorStmt = &Interpreter{}

func (i *Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		err := i.execute(statement)
		// the return value of execute is either nil or a runtime error
		if err != nil {
			fmt.Println(err.(*runtimeError).Error())
			return
		}
	}
}

func (i *Interpreter) execute(stmt Stmt) interface{} {
	return stmt.Accept(i)
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) visitExpressionStmt(stmt *ExpressionStmt) interface{} {
	expr := i.evaluate(stmt.expression)
	err, ok := expr.(*runtimeError)
	if ok {
		return err
	}
	return nil
}

func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) interface{} {
	value := i.evaluate(stmt.expression)
	err, ok := value.(*runtimeError)
	if ok {
		return err
	}
	fmt.Println(value)
	return nil
}

func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) interface{} {
	left := i.evaluate(expr.left)
	err, ok := left.(*runtimeError)
	if ok {
		return err
	}
	right := i.evaluate(expr.right)
	err, ok = right.(*runtimeError)
	if ok {
		return err
	}

	switch expr.operator.Type {
	case EqualEqual:
		return isEqual(left, right)
	case BangEqual:
		return !isEqual(left, right)
	case Greater:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft > castedRight
	case GreaterEqual:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft >= castedRight
	case Less:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft > castedRight
	case LessEqual:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft < castedRight
	case Slash:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft / castedRight
	case Star:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft * castedRight
	case Minus:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft - castedRight
	case Plus:
		leftFloat64, rightFloat64, err := castNumberOperands(expr.operator, left, right)
		if err == nil {
			return leftFloat64 + rightFloat64
		}
		leftString, rightString, err := castStringOperands(expr.operator, left, right)
		if err == nil {
			return leftString + rightString
		}
		return &runtimeError{token: expr.operator, message: "Operands must be two numbers or two strings."}
	}

	// unreachable
	return nil
}

func (i *Interpreter) visitGroupingExpr(expr *GroupingExpr) interface{} {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) interface{} {
	return expr.value
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) interface{} {
	right := i.evaluate(expr.right)
	err, ok := right.(*runtimeError)
	if ok {
		return err
	}

	switch expr.operator.Type {
	case Bang:
		return !isTruthy(right)
	case Minus:
		castedRight, err := castNumberOperand(expr.operator, right)
		if err != nil {
			return err
		}
		return -castedRight
	}

	// unreachable
	return nil
}

func isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}
	if value == false {
		return false
	}
	return true
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a == b
}

type runtimeError struct {
	token   Token
	message string
}

func (e *runtimeError) Error() string {
	return fmt.Sprintf("%s: %s\n[line %d]", e.token.Lexeme, e.message, e.token.Line)
}

func castNumberOperand(operator Token, operand interface{}) (float64, *runtimeError) {
	casted, ok := operand.(float64)
	if !ok {
		return 0, &runtimeError{token: operator, message: "Operand must be a number."}
	}
	return casted, nil
}

func castNumberOperands(operator Token, left, right interface{}) (float64, float64, *runtimeError) {
	castedLeft, okLeft := left.(float64)
	castedRight, okRight := right.(float64)
	if !okLeft || !okRight {
		return 0, 0, &runtimeError{token: operator, message: "Operands must be numbers."}
	}
	return castedLeft, castedRight, nil
}

func castStringOperands(operator Token, left, right interface{}) (string, string, *runtimeError) {
	castedLeft, okLeft := left.(string)
	castedRight, okRight := right.(string)
	if !okLeft || !okRight {
		return "", "", &runtimeError{token: operator, message: "Operands must be strings."}
	}
	return castedLeft, castedRight, nil
}
