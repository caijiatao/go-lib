package compilation_principle

import (
	"go/ast"
	"go/parser"
)

func Parser() {
	expr, _ := parser.ParseExpr(`1+2*3`)

	ast.Print(nil, expr)
	expr, _ = parser.ParseExpr(`x`)
	ast.Print(nil, expr)
}
