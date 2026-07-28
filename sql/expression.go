package sql

import (
	"fmt"
	"strings"
)

type SqlExpression interface {
	String() string
}

// ValueExpression value
type ValueExpression struct {
	Value any
}

func (v *ValueExpression) String() string {
	return fmt.Sprintf("%v", v.Value)
}

// ReferenceExpression correlation.columname
type ReferenceExpression struct {
	Correlation string
	Name        string
}

func (p *ReferenceExpression) String() string {
	if len(p.Correlation) <= 0 {
		return p.Name
	}

	return fmt.Sprintf("%s.%s", p.Correlation, p.Name)
}

// PositionalExpression $number
type PositionalExpression struct {
	Position int
}

func (p *PositionalExpression) String() string {
	return fmt.Sprintf("$%d", p.Position)
}

// SubscriptExpression expression[subscript]
type SubscriptExpression struct {
	Expression SqlExpression
	Subscript  int
}

func (p *SubscriptExpression) String() string {
	return fmt.Sprintf("%s[%d]", p.Expression, p.Subscript)
}

// SubscriptSliceExpression expression[lower_subscript:upper_subscript]
type SubscriptSliceExpression struct {
	Expression     SqlExpression
	LowerSubscript int
	UpperSubscript int
}

func (p *SubscriptSliceExpression) String() string {
	return fmt.Sprintf("%s[%d:%d]", p.Expression, p.LowerSubscript, p.UpperSubscript)
}

// FunctionExpression function(expression, ...expression)
type FunctionExpression struct {
	FunctionName string
	Parameters   []SqlExpression
}

func (s *FunctionExpression) String() string {
	items := make([]string, len(s.Parameters))
	for i, parameter := range s.Parameters {
		items[i] = parameter.String()
	}

	return fmt.Sprintf("%s(%s)", s.FunctionName, strings.Join(items, ", "))
}

// SubExpression (expression)
type SubExpression struct {
	Expression SqlExpression
}

func (s *SubExpression) String() string {
	return fmt.Sprintf("(%s)", s.Expression)
}

// AliasExpression expression AS alias OR expression AS "alias"
type AliasExpression struct {
	Expression SqlExpression
	Alias      string
	Quoted     bool
}

func (s *AliasExpression) String() string {
	if s.Quoted {
		return fmt.Sprintf(`%s AS "%s"`, s.Expression, s.Alias)
	}

	return fmt.Sprintf("%s AS %s", s.Expression, s.Alias)
}

// TypeAliasExpression expression AS alias OR expression AS typename
type TypeAliasExpression struct {
	Expression SqlExpression
	TypeName   string
}

func (s *TypeAliasExpression) String() string {
	return fmt.Sprintf("%s AS %s", s.Expression, s.TypeName)
}

// FieldExpression expression.fieldname
type FieldExpression struct {
	Expression SqlExpression
	FieldName  string
}

func (s *FieldExpression) String() string {
	return fmt.Sprintf("%s.%s", s.Expression, s.FieldName)
}

// QueryExpression subquery
type QueryExpression struct {
	Statement *SqlStatement
}

func (s *QueryExpression) String() string {
	return s.Statement.String()
}

// WildcardExpression * OR correlation.*
type WildcardExpression struct {
	Correlation string
}

func (s *WildcardExpression) String() string {
	if len(s.Correlation) > 0 {
		return fmt.Sprintf("%s.*", s.Correlation)
	}

	return "*"
}

// CombinatoryExpression expression <combinatory> expression
type CombinatoryExpression struct {
	Left     SqlExpression
	Right    SqlExpression
	Operator string
}

func (s *CombinatoryExpression) String() string {
	return fmt.Sprintf("%s %s %s", s.Left, s.Operator, s.Right)
}
