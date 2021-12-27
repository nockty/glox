package lox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinter(t *testing.T) {
	testCases := []struct {
		expr     Expr
		expected string
	}{
		{
			expr: NewBinaryExpr(
				NewUnaryExpr(
					NewToken(Minus, "-", nil, 1),
					NewLiteralExpr(123),
				),
				NewToken(Star, "*", nil, 1),
				NewGroupingExpr(NewLiteralExpr(45.67)),
			),
			expected: "(* (- 123) (group 45.67))",
		},
		{
			expr: NewBinaryExpr(
				NewBinaryExpr(
					NewLiteralExpr(42),
					NewToken(Plus, "+", nil, 1),
					NewBinaryExpr(
						NewLiteralExpr(50),
						NewToken(Star, "*", nil, 1),
						NewGroupingExpr(
							NewBinaryExpr(
								NewLiteralExpr(1),
								NewToken(Plus, "+", nil, 1),
								NewLiteralExpr(5),
							),
						),
					),
				),
				NewToken(Minus, "-", nil, 1),
				NewBinaryExpr(
					NewLiteralExpr(9),
					NewToken(Slash, "/", nil, 1),
					NewLiteralExpr(3),
				),
			),
			expected: "(- (+ 42 (* 50 (group (+ 1 5)))) (/ 9 3))",
		},
	}

	for _, tc := range testCases {
		actual := (&AstPrinter{}).Sprint(tc.expr)
		assert.Equal(t, tc.expected, actual)
	}
}
