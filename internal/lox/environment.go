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

func (e *environment) assign(name Token, value interface{}) *runtimeError {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = value
		return nil
	}
	return &runtimeError{
		token:   name,
		message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
	}
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
