package parser

import (
	"interpreter/ast"
	"interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	lexer := lexer.New(input)
	p := New(lexer)
	program := p.ParseProgram()
	checkParserError(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() return nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}

}

func TestReturnStatements(t *testing.T) {
	input := `
			return 5;
			return 4;
			return 123134123;
		`
	lex := lexer.New(input)
	parser := New(lex)
	checkParserError(t, parser)
	program := parser.ParseProgram()

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}

}

func checkParserError(t *testing.T, p *Parser) {
	if len(p.errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.errors))
	for _, value := range p.errors {
		t.Errorf("parser error: %q", value)

	}
	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmt.Name.TokenLiteral())
		return false
	}
	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := `
	myVar;
	`
	l := lexer.New(input)
	p := New(l)
	checkParserError(t, p)

	pr := p.ParseProgram()

	if len(pr.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(pr.Statements))
	}

	stmt, ok := pr.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", pr.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.TokenLiteral() != "myVar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "myVar",
			ident.TokenLiteral())
	}

	if ident.Value != "myVar" {
		t.Errorf("ident.Value not %s. got=%s", "myVar", ident.Value)
	}

}
