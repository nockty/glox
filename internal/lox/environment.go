package lox

import "fmt"

type environment struct {
	values map[string]interface{}
}

func newEnvironment() *environment {
	return &environment{values: make(map[string]interface{})}
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *environment) get(name Token) (interface{}, *runtimeError) {
	value, ok := e.values[name.Lexeme]
	if !ok {
		return nil, &runtimeError{
			token:   name,
			message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
		}
	}
	return value, nil
}
