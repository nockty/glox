package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstPrinter(t *testing.T) {
	expr := &Binary{
		NewUnary(
			NewToken(Minus, "-", nil, 1),
			NewLiteral(123),
		),
		NewToken(Star, "*", nil, 1),
		NewGrouping(NewLiteral(45.67)),
	}
	expected := "(* (- 123) (group 45.67))"
	actual := (&astPrinter{}).Print(expr)
	assert.Equal(t, expected, actual)
}
