package lox

import "fmt"

type interpreter struct {
	env *environment
}

// interpreter implements visitorExpr and visitorStmt
var _ visitorExpr = &interpreter{}
var _ visitorStmt = &interpreter{}

func NewInterpreter() *interpreter {
	return &interpreter{env: newEnvironment()}
}

func (i *interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			fmt.Println(err.(*runtimeError).Error())
			return
		}
	}
}

// execute returns either nil or a runtime error
func (i *interpreter) execute(stmt Stmt) interface{} {
	return stmt.Accept(i)
}

func (i *interpreter) executeBlock(statements []Stmt, env *environment) interface{} {
	previous := i.env
	defer func() { i.env = previous }()
	i.env = env
	for _, statement := range statements {
		err := i.execute(statement)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) visitBlockStmt(stmt *BlockStmt) interface{} {
	return i.executeBlock(stmt.statements, newScopedEnvironment(i.env))
}

func (i *interpreter) visitExpressionStmt(stmt *ExpressionStmt) interface{} {
	expr := i.evaluate(stmt.expression)
	err, ok := expr.(*runtimeError)
	if ok {
		return err
	}
	return nil
}

func (i *interpreter) visitIfStmt(stmt *IfStmt) interface{} {
	condition := i.evaluate(stmt.condition)
	err, ok := condition.(*runtimeError)
	if ok {
		return err
	}
	if isTruthy(condition) {
		err := i.execute(stmt.thenBranch)
		if err != nil {
			return err
		}
	} else if stmt.elseBranch != nil {
		err := i.execute(stmt.thenBranch)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) visitPrintStmt(stmt *PrintStmt) interface{} {
	value := i.evaluate(stmt.expression)
	err, ok := value.(*runtimeError)
	if ok {
		return err
	}
	fmt.Println(stringify(value))
	return nil
}

func (i *interpreter) visitVarStmt(stmt *VarStmt) interface{} {
	var value interface{} = nil
	if stmt.initializer != nil {
		value = i.evaluate(stmt.initializer)
	}
	i.env.define(stmt.name.Lexeme, value)
	return nil
}

func (i *interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *interpreter) visitAssignExpr(expr *AssignExpr) interface{} {
	value := i.evaluate(expr.value)
	err, ok := value.(*runtimeError)
	if ok {
		return err
	}
	err = i.env.assign(expr.name, value)
	if err != nil {
		return err
	}
	return value
}

func (i *interpreter) visitBinaryExpr(expr *BinaryExpr) interface{} {
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
		return castedLeft < castedRight
	case LessEqual:
		castedLeft, castedRight, err := castNumberOperands(expr.operator, left, right)
		if err != nil {
			return err
		}
		return castedLeft <= castedRight
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

func (i *interpreter) visitGroupingExpr(expr *GroupingExpr) interface{} {
	return i.evaluate(expr.expression)
}

func (i *interpreter) visitLiteralExpr(expr *LiteralExpr) interface{} {
	return expr.value
}

func (i *interpreter) visitLogicalExpr(expr *LogicalExpr) interface{} {
	left := i.evaluate(expr.left)
	err, ok := left.(*runtimeError)
	if ok {
		return err
	}
	if expr.operator.Type == Or {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.right)
}

func (i *interpreter) visitUnaryExpr(expr *UnaryExpr) interface{} {
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

func (i *interpreter) visitVariableExpr(expr *VariableExpr) interface{} {
	value, err := i.env.get(expr.name)
	if err != nil {
		return err
	}
	return value
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

func stringify(value interface{}) string {
	if value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", value)
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
