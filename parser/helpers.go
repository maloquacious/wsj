// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package parser

// helper functions for the generated parser

import (
	"github.com/maloquacious/wsj/ast"
)

// bdup returns a copy of a slice
func bdup(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// toAnySlice helps us navigate Pigeon's nodes
func toAnySlice(v any) []any {
	if v == nil {
		return nil
	}
	return v.([]any)
}

// coerceExprList ensures that the expression list is never nil
func coerceExprList(v any) []ast.Expr {
	if v == nil {
		return []ast.Expr{}
	}
	return v.([]ast.Expr)
}

func coerceStatementList(v any) []ast.Stmt {
	if v == nil {
		return []ast.Stmt{}
	}
	items := v.([]interface{})
	list := make([]ast.Stmt, len(items))
	for i, s := range items {
		list[i] = s.(ast.Stmt)
	}
	return list
}

func coerceSuffixList(v any) []ast.Suffix {
	if v == nil {
		return []ast.Suffix{}
	}
	items := v.([]interface{})
	list := make([]ast.Suffix, len(items))
	for i, s := range items {
		list[i] = s.(ast.Suffix)
	}
	return list
}

func foldLeftBinary(left any, rest []interface{}) (ast.Expr, error) {
	result := left.(ast.Expr)

	for _, part := range rest {
		pair := part.([]interface{})
		// pair[0] is the operator match (like (_ "+" _)), pair[1] is the right operand
		opMatch := pair[0].([]interface{})
		// opMatch[1] is the actual operator string (surrounded by whitespace)
		op := string(opMatch[1].([]byte))
		right := pair[1].(ast.Expr)

		result = &ast.BinaryExpr{
			Left:     result,
			Operator: op,
			Right:    right,
			Pos:      result.Position(), // Or use position of operator if available
		}
	}

	return result, nil
}

func foldRightBinary(left any, rest []interface{}) (ast.Expr, error) {
	// Base case: no operators
	if len(rest) == 0 {
		return left.(ast.Expr), nil
	}

	// Right-to-left fold
	// Start from the rightmost pair and work backward
	last := rest[len(rest)-1].([]interface{})
	op := last[0].(string)
	right := last[1].(ast.Expr)

	// Fold tail recursively
	headExpr, err := foldRightBinary(left, rest[:len(rest)-1])
	if err != nil {
		return nil, err
	}

	return &ast.BinaryExpr{
		Left:     headExpr,
		Operator: op,
		Right:    right,
		Pos:      headExpr.Position(), // Optional: use operator position if available
	}, nil
}
