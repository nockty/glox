package lox

import "fmt"

type Interpreter struct{}

// Interpreter implements visitor
var _ visitorExpr = &Interpreter{}

func (i *Interpreter) Interpret(expr Expr) {
	value := i.evaluate(expr)
	err, ok := value.(*runtimeError)
	if ok {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(value)
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
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
