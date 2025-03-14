package main

import (
	"testing"

	"github.com/MKode312/calculator_go/pkg/ast"
)

func TestParseAST_ValidExpression(t *testing.T) {
	expressions := []struct {
		input    string
		expected *ast.ASTNode
	}{
		{
			input: "1 + 2",
			expected: &ast.ASTNode{
				IsLeaf:   false,
				Operator: "+",
				Left: &ast.ASTNode{
					IsLeaf: true,
					Value:  1,
				},
				Right: &ast.ASTNode{
					IsLeaf: true,
					Value:  2,
				},
			},
		},
		{
			input: "3 - 4 * 5",
			expected: &ast.ASTNode{
				IsLeaf:   false,
				Operator: "-",
				Left: &ast.ASTNode{
					IsLeaf: true,
					Value:  3,
				},
				Right: &ast.ASTNode{
					IsLeaf:   false,
					Operator: "*",
					Left: &ast.ASTNode{
						IsLeaf: true,
						Value:  4,
					},
					Right: &ast.ASTNode{
						IsLeaf: true,
						Value:  5,
					},
				},
			},
		},
		{
			input: "(1 + 2) * 3",
			expected: &ast.ASTNode{
				IsLeaf:   false,
				Operator: "*",
				Left: &ast.ASTNode{
					IsLeaf:   false,
					Operator: "+",
					Left: &ast.ASTNode{
						IsLeaf: true,
						Value:  1,
					},
					Right: &ast.ASTNode{
						IsLeaf: true,
						Value:  2,
					},
				},
				Right: &ast.ASTNode{
					IsLeaf: true,
					Value:  3,
				},
			},
		},
	}

	for _, expr := range expressions {
		result, err := ast.ParseAST(expr.input)
		if err != nil {
			t.Errorf("ParseAST(%q) returned error: %v", expr.input, err)
			continue
		}
		if !compareASTNodes(result, expr.expected) {
			t.Errorf("ParseAST(%q) = %v, expected %v", expr.input, result, expr.expected)
		}
	}
}

func TestParseAST_InvalidExpression(t *testing.T) {
	invalidExpressions := []struct {
		input string
	}{
		{" "},
		{"1 +"},
		{"(1 + 2"},
		{"1 + (2 * 3))"},
	}

	for _, expr := range invalidExpressions {
		_, err := ast.ParseAST(expr.input)
		if err == nil {
			t.Errorf("ParseAST(%q) expected an error but got nil", expr.input)
		}
	}
}

func compareASTNodes(a, b *ast.ASTNode) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if a.IsLeaf != b.IsLeaf || a.Operator != b.Operator || a.Value != b.Value {
		return false
	}
	return compareASTNodes(a.Left, b.Left) && compareASTNodes(a.Right, b.Right)
}
