package parser

import (
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/stretchr/testify/assert"
)

func TestParseExprNoErrors(t *testing.T) {
	tests := map[string]struct {
		input string
	}{
		"StringLiteral": {
			input: "\"hello\"",
		},
		"NumberLiteral": {
			input: "5",
		},
		"Addition": {
			input: "a + b",
		},
		"AddSub": {
			input: "a - b + c",
		},
		"MulAdd": {
			input: "a * b + c * d",
		},
		"MulDiv": {
			input: "a / b * c",
		},
		"UnaryOps": {
			input: "+a - -b",
		},
		"Parens": {
			input: "a * (b + c)",
		},
		"Call": {
			input: "foo(a, b, c)",
		},
		"CallPrecedence": {
			input: "a + foo(b)",
		},
		"CurriedCall": {
			input: "foo(a)(b)(c)",
		},
		"OptChainCall": {
			input: "foo?(bar)",
		},
		"ArrayLiteral": {
			input: "[1, 2, 3]",
		},
		"Member": {
			input: "a.b?.c",
		},
		"MemberPrecedence": {
			input: "a + b.c",
		},
		"Index": {
			input: "a[base + offset]",
		},
		"IndexPrecedence": {
			input: "a + b[c]",
		},
		"MultipleIndexes": {
			input: "a[i][j]",
		},
		"OptChainIndex": {
			input: "a?[base + offset]",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			source := Source{
				path:     "input.esc",
				Contents: test.input,
			}

			parser := NewParser(source)
			expr, _ := parser.parseExpr()

			snaps.MatchSnapshot(t, expr)
			assert.Len(t, parser.errors, 0)
		})
	}
}

func TestParseExprErrorHandling(t *testing.T) {
	tests := map[string]struct {
		input string
	}{
		"IncompleteBinaryExpr": {
			input: "a - b +",
		},
		"ExtraOperatorsInBinaryExpr": {
			input: "a + * b",
		},
		"IncompleteCall": {
			input: "foo(a,",
		},
		"IncompleteMember": {
			input: "a + b.",
		},
		"IncompleteMemberOptChain": {
			input: "a + b?.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			source := Source{
				path:     "input.esc",
				Contents: test.input,
			}

			parser := NewParser(source)
			expr, _ := parser.parseExpr()

			snaps.MatchSnapshot(t, expr)
			assert.Len(t, parser.errors, 1)
			snaps.MatchSnapshot(t, parser.errors[0])
		})
	}
}
