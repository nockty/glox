package lox

import "fmt"

type environment struct {
	enclosing *environment
	values    map[string]interface{}
}

func newEnvironment() *environment {
	return &environment{
		enclosing: nil,
		values:    make(map[string]interface{}),
	}
}

func newScopedEnvironment(enclosing *environment) *environment {
	return &environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *environment) define(name string, value interface{}) {
	e.values[name] = value
}

func (e *environment) assign(name Token, value interface{}) *runtimeError {
	if _, ok := e.values[name.Lexeme]; !ok {
		if e.enclosing != nil {
			return e.enclosing.assign(name, value)
		}
		return &runtimeError{
			token:   name,
			message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
		}
	}
	e.values[name.Lexeme] = value
	return nil
}

func (e *environment) get(name Token) (interface{}, *runtimeError) {
	value, ok := e.values[name.Lexeme]
	if !ok {
		if e.enclosing != nil {
			return e.enclosing.get(name)
		}
		return nil, &runtimeError{
			token:   name,
			message: fmt.Sprintf("Undefined variable '%s'.", name.Lexeme),
		}
	}
	return value, nil
}
