package parser

import (
	"testing"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/lexer"
	"github.com/MahmoudHilani/GoInterpreter/test/src/monkey/ast"
)

func testLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l:= lexer.New(input)
	p:= New(l)

	program := p.ParseProgram()

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements doesn't contain 3 components, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"}, {"y"}, {"foobar"},
	}

	for i, tt := range(tests) {
		stmt :=  program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	} 

}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not let, got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok:= s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatemant. Got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got '%s' instead", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s' got '%s' instead", name, letStmt.Name)
		return false
	}

	return true
}