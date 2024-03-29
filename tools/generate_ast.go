package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path"
	"strings"
)

var goTypes = [...]string{
	"bool",
	"string",
	"int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
	"byte", // alias for uint8
	"rune", // alias for int32
	"float32", "float64",
	"complex64", "complex128",
}

func main() {
	args := os.Args
	if len(args) != 2 {
		println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outDir := args[1]
	types := []string{
		"Assign   : name Token, value Expr",
		"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Logical  : left Expr, operator Token, right Expr",
		"Unary    : operator Token, right Expr",
		"Variable : name Token",
	}
	err := defineAST(outDir, "Expr", types)
	if err != nil {
		println(fmt.Errorf("failed to generate Expr AST: %w", err).Error())
	}
	types = []string{
		"Block      : statements []Stmt",
		"Expression : expression Expr",
		"If         : condition Expr, thenBranch Stmt, elseBranch Stmt",
		"Print      : expression Expr",
		"Var        : name Token, initializer Expr",
		"While      : condition Expr, body Stmt",
	}
	err = defineAST(outDir, "Stmt", types)
	if err != nil {
		println(fmt.Errorf("failed to generate Stmt AST: %w", err).Error())
	}
}

func writeLines(w *bufio.Writer, lines []string) error {
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return fmt.Errorf("write line: %w", err)
		}
	}
	return w.Flush()
}

func defineAST(outDir, baseName string, types []string) error {
	path := path.Join(outDir, strings.ToLower(baseName)+".go")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %v: %w", path, err)
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buf)
	lines := []string{
		// This comment should match the regex in https://github.com/golang/go/issues/13560#issuecomment-288457920
		"// Code generated by the generate_ast tool; DO NOT EDIT.",
		"package lox",
		"",
	}
	err = writeLines(w, lines)
	if err != nil {
		return err
	}
	err = defineBase(w, baseName)
	if err != nil {
		return fmt.Errorf("define Expr: %w", err)
	}
	err = defineVisitor(w, baseName, types)
	if err != nil {
		return fmt.Errorf("define visitor: %w", err)
	}
	// types
	for _, exprType := range types {
		components := strings.Split(exprType, ":")
		typeName := strings.TrimSpace(components[0]) + strings.Title(baseName)
		fields := strings.TrimSpace(components[1])
		err := defineType(w, baseName, typeName, fields)
		if err != nil {
			return fmt.Errorf("define type %v: %w", typeName, err)
		}
	}
	// format source code
	source, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format code: %w", err)
	}
	_, err = f.Write(source)
	if err != nil {
		return fmt.Errorf("write code to file %v: %w", f, err)
	}
	return nil
}

func defineBase(w *bufio.Writer, baseName string) error {
	lines := []string{
		fmt.Sprintf("type %s interface {", strings.Title(baseName)),
		fmt.Sprintf("Accept(v visitor%s) interface{}", strings.Title(baseName)),
	}
	// Accept by go type
	for _, goType := range goTypes {
		lines = append(lines, fmt.Sprintf("Accept%s(v visitor%s%s) %s",
			strings.Title(goType), strings.Title(baseName), strings.Title(goType), goType))
	}
	lines = append(lines, "}", "")
	return writeLines(w, lines)
}

func defineVisitor(w *bufio.Writer, baseName string, types []string) error {
	lines := []string{fmt.Sprintf("type visitor%s interface {", strings.Title(baseName))}
	for _, exprType := range types {
		typeName := strings.TrimSpace(strings.Split(exprType, ":")[0]) + strings.Title(baseName)
		lines = append(lines, fmt.Sprintf(
			"visit%s(*%s) interface{}", typeName, typeName))
	}
	lines = append(lines, "}", "")
	// visitor by go type
	for _, goType := range goTypes {
		lines = append(lines, fmt.Sprintf("type visitor%s%s interface {", strings.Title(baseName), strings.Title(goType)))
		for _, exprType := range types {
			typeName := strings.TrimSpace(strings.Split(exprType, ":")[0]) + strings.Title(baseName)
			lines = append(lines, fmt.Sprintf(
				"visit%s(*%s) %s", typeName, typeName, goType))
		}
		lines = append(lines, "}", "")
	}
	return writeLines(w, lines)
}

func defineType(w *bufio.Writer, baseName, typeName, fieldList string) error {
	fieldsUntrimmed := strings.Split(fieldList, ",")
	fields := []string{}
	for _, field := range fieldsUntrimmed {
		fields = append(fields, strings.TrimSpace(field))
	}
	// type
	lines := []string{
		fmt.Sprintf("type %s struct {", typeName),
	}
	lines = append(lines, fields...)
	lines = append(lines, "}", "")
	// implements Expr
	lines = append(lines,
		fmt.Sprintf("// %s implements %s", typeName, strings.Title(baseName)),
		fmt.Sprintf("var _ %s = &%s{}", strings.Title(baseName), typeName),
		"",
	)
	// constructor
	lines = append(lines,
		fmt.Sprintf("func New%s(%s) *%s {", typeName, fieldList, typeName),
		fmt.Sprintf("return &%s {", typeName),
	)
	for _, field := range fields {
		name := strings.Split(field, " ")[0]
		lines = append(lines, fmt.Sprintf("%s: %s,", name, name))
	}
	lines = append(lines, "}", "}", "")
	// visitor pattern
	lines = append(lines,
		fmt.Sprintf("func (expr *%s) Accept(v visitor%s) interface{} {", typeName, strings.Title(baseName)),
		fmt.Sprintf("return v.visit%s(expr)", typeName),
		"}",
		"",
	)
	// visitor pattern by go type
	for _, goType := range goTypes {
		lines = append(lines,
			fmt.Sprintf("func (expr *%s) Accept%s(v visitor%s%s) %s {",
				typeName, strings.Title(goType), strings.Title(baseName), strings.Title(goType), goType),
			fmt.Sprintf("return v.visit%s(expr)", typeName),
			"}",
			"",
		)
	}
	return writeLines(w, lines)
}
