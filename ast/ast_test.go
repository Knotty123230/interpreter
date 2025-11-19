package ast

import (
	"interpreter/token"
	"testing"
)

func TestAstString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "hello",
					},
					Value: "hello",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "myVar",
					},
					Value: "myVar",
				},
			},
		},
	}
	if program.String() != "let hello = myVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
