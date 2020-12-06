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

func main() {
	args := os.Args
	if len(args) != 2 {
		println("Usage: generate_ast <output directory>")
		os.Exit(64)
	}
	outDir := args[1]
	types := []string{
		"Binary   : left Expr, operator Token, right Expr",
		"Grouping : expression Expr",
		"Literal  : value interface{}",
		"Unary    : operator Token, right Expr",
	}
	err := defineAST(outDir, "expressions", types)
	if err != nil {
		println(fmt.Errorf("failed to generate AST: %v", err).Error())
	}
}

func writeLines(w *bufio.Writer, lines []string) error {
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			return fmt.Errorf("write line: %v", err)
		}
	}
	return w.Flush()
}

func defineAST(outDir, baseName string, types []string) error {
	path := path.Join(outDir, baseName+".go")
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %v: %v", path, err)
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	w := bufio.NewWriter(buf)
	lines := []string{
		// This comment should match the regex in https://github.com/golang/go/issues/13560#issuecomment-288457920
		"// Code generated by the generate_ast tool; DO NOT EDIT.",
		"package glox",
		"",
		"type Expr interface {",
		"accept(v Visitor) interface{}",
		"}",
		"",
	}
	err = writeLines(w, lines)
	if err != nil {
		return err
	}
	err = defineVisitor(w, types)
	if err != nil {
		return fmt.Errorf("define visitor: %v", err)
	}
	// types
	for _, exprType := range types {
		components := strings.Split(exprType, ":")
		typeName := strings.TrimSpace(components[0])
		fields := strings.TrimSpace(components[1])
		err := defineType(w, typeName, fields)
		if err != nil {
			return fmt.Errorf("define type %v: %v", typeName, err)
		}
	}
	// format source code
	source, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format code: %v", err)
	}
	_, err = f.Write(source)
	if err != nil {
		return fmt.Errorf("write code to file %v: %v", f, err)
	}
	return nil
}

func defineVisitor(w *bufio.Writer, types []string) error {
	lines := []string{"type Visitor interface {"}
	for _, exprType := range types {
		typeName := strings.TrimSpace(strings.Split(exprType, ":")[0])
		lines = append(lines, fmt.Sprintf("visit%s(*%s) interface{}", typeName, typeName))
	}
	lines = append(lines, "}", "")
	return writeLines(w, lines)
}

func defineType(w *bufio.Writer, typeName, fieldList string) error {
	fieldsUntrimmed := strings.Split(fieldList, ",")
	fields := []string{}
	for _, field := range fieldsUntrimmed {
		fields = append(fields, strings.TrimSpace(field))
	}
	// type
	lines := []string{
		fmt.Sprintf("type %s struct {", typeName),
	}
	for _, field := range fields {
		lines = append(lines, field)
	}
	lines = append(lines, "}", "")
	// implements Expr
	lines = append(lines,
		fmt.Sprintf("// %s implements Expr", typeName),
		fmt.Sprintf("var _ Expr = &%s{}", typeName),
		"",
	)
	// visitor pattern
	lines = append(lines,
		fmt.Sprintf("func (expr *%s) accept(v Visitor) interface{} {", typeName),
		fmt.Sprintf("return v.visit%s(expr)", typeName),
		"}",
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
	return writeLines(w, lines)
}
