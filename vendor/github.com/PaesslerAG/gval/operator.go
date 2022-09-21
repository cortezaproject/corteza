package gval

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

type stage struct {
	Evaluable
	infixBuilder
	operatorPrecedence
}

type stageStack []stage //operatorPrecedence in stacktStage is continuously, monotone ascending

func (s *stageStack) push(b stage) error {
	for len(*s) > 0 && s.peek().operatorPrecedence >= b.operatorPrecedence {
		a := s.pop()
		eval, err := a.infixBuilder(a.Evaluable, b.Evaluable)
		if err != nil {
			return err
		}
		if a.IsConst() && b.IsConst() {
			v, err := eval(nil, nil)
			if err != nil {
				return err
			}
			b.Evaluable = constant(v)
			continue
		}
		b.Evaluable = eval
	}
	*s = append(*s, b)
	return nil
}

func (s *stageStack) peek() stage {
	return (*s)[len(*s)-1]
}

func (s *stageStack) pop() stage {
	a := s.peek()
	(*s) = (*s)[:len(*s)-1]
	return a
}

type infixBuilder func(a, b Evaluable) (Evaluable, error)

func (l Language) isSymbolOperation(r rune) bool {
	_, in := l.operatorSymbols[r]
	return in
}

func (l Language) isOperatorPrefix(op string) bool {
	for k := range l.operators {
		if strings.HasPrefix(k, op) {
			return true
		}
	}
	return false
}

func (op *infix) initiate(name string) {
	f := func(a, b interface{}) (interface{}, error) {
		return nil, fmt.Errorf("invalid operation (%T) %s (%T)", a, name, b)
	}
	if op.arbitrary != nil {
		f = op.arbitrary
	}
	for _, typeConvertion := range []bool{true, false} {
		if op.text != nil && (!typeConvertion || op.arbitrary == nil) {
			f = getStringOpFunc(op.text, f, typeConvertion)
		}
		if op.boolean != nil {
			f = getBoolOpFunc(op.boolean, f, typeConvertion)
		}
		if op.number != nil {
			f = getFloatOpFunc(op.number, f, typeConvertion)
		}
		if op.decimal != nil {
			f = getDecimalOpFunc(op.decimal, f, typeConvertion)
		}
	}
	if op.shortCircuit == nil {
		op.builder = func(a, b Evaluable) (Evaluable, error) {
			return func(c context.Context, x interface{}) (interface{}, error) {
				a, err := a(c, x)
				if err != nil {
					return nil, err
				}
				b, err := b(c, x)
				if err != nil {
					return nil, err
				}
				return f(a, b)
			}, nil
		}
		return
	}
	shortF := op.shortCircuit
	op.builder = func(a, b Evaluable) (Evaluable, error) {
		return func(c context.Context, x interface{}) (interface{}, error) {
			a, err := a(c, x)
			if err != nil {
				return nil, err
			}
			if r, ok := shortF(a); ok {
				return r, nil
			}
			b, err := b(c, x)
			if err != nil {
				return nil, err
			}
			return f(a, b)
		}, nil
	}
}

type opFunc func(a, b interface{}) (interface{}, error)

func getStringOpFunc(s func(a, b string) (interface{}, error), f opFunc, typeConversion bool) opFunc {
	if typeConversion {
		return func(a, b interface{}) (interface{}, error) {
			if a != nil && b != nil {
				return s(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
			}
			return f(a, b)
		}
	}
	return func(a, b interface{}) (interface{}, error) {
		s1, k := a.(string)
		s2, l := b.(string)
		if k && l {
			return s(s1, s2)
		}
		return f(a, b)
	}
}
func convertToBool(o interface{}) (bool, bool) {
	if b, ok := o.(bool); ok {
		return b, true
	}
	v := reflect.ValueOf(o)
	for o != nil && v.Kind() == reflect.Ptr {
		v = v.Elem()
		if !v.IsValid() {
			return false, false
		}
		o = v.Interface()
	}
	if o == false || o == nil || o == "false" || o == "FALSE" {
		return false, true
	}
	if o == true || o == "true" || o == "TRUE" {
		return true, true
	}
	if f, ok := convertToFloat(o); ok {
		return f != 0., true
	}
	return false, false
}
func getBoolOpFunc(o func(a, b bool) (interface{}, error), f opFunc, typeConversion bool) opFunc {
	if typeConversion {
		return func(a, b interface{}) (interface{}, error) {
			x, k := convertToBool(a)
			y, l := convertToBool(b)
			if k && l {
				return o(x, y)
			}
			return f(a, b)
		}
	}
	return func(a, b interface{}) (interface{}, error) {
		x, k := a.(bool)
		y, l := b.(bool)
		if k && l {
			return o(x, y)
		}
		return f(a, b)
	}
}
func convertToFloat(o interface{}) (float64, bool) {
	if i, ok := o.(float64); ok {
		return i, true
	}
	v := reflect.ValueOf(o)
	for o != nil && v.Kind() == reflect.Ptr {
		v = v.Elem()
		if !v.IsValid() {
			return 0, false
		}
		o = v.Interface()
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	}
	if s, ok := o.(string); ok {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return f, true
		}
	}
	return 0, false
}
func getFloatOpFunc(o func(a, b float64) (interface{}, error), f opFunc, typeConversion bool) opFunc {
	if typeConversion {
		return func(a, b interface{}) (interface{}, error) {
			x, k := convertToFloat(a)
			y, l := convertToFloat(b)
			if k && l {
				return o(x, y)
			}

			return f(a, b)
		}
	}
	return func(a, b interface{}) (interface{}, error) {
		x, k := a.(float64)
		y, l := b.(float64)
		if k && l {
			return o(x, y)
		}

		return f(a, b)
	}
}
func convertToDecimal(o interface{}) (decimal.Decimal, bool) {
	if i, ok := o.(decimal.Decimal); ok {
		return i, true
	}
	if i, ok := o.(float64); ok {
		return decimal.NewFromFloat(i), true
	}
	v := reflect.ValueOf(o)
	for o != nil && v.Kind() == reflect.Ptr {
		v = v.Elem()
		if !v.IsValid() {
			return decimal.Zero, false
		}
		o = v.Interface()
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return decimal.NewFromInt(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return decimal.NewFromFloat(float64(v.Uint())), true
	case reflect.Float32, reflect.Float64:
		return decimal.NewFromFloat(v.Float()), true
	}
	if s, ok := o.(string); ok {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			return decimal.NewFromFloat(f), true
		}
	}
	return decimal.Zero, false
}
func getDecimalOpFunc(o func(a, b decimal.Decimal) (interface{}, error), f opFunc, typeConversion bool) opFunc {
	if typeConversion {
		return func(a, b interface{}) (interface{}, error) {
			x, k := convertToDecimal(a)
			y, l := convertToDecimal(b)
			if k && l {
				return o(x, y)
			}

			return f(a, b)
		}
	}
	return func(a, b interface{}) (interface{}, error) {
		x, k := a.(decimal.Decimal)
		y, l := b.(decimal.Decimal)
		if k && l {
			return o(x, y)
		}

		return f(a, b)
	}
}

type operator interface {
	merge(operator) operator
	precedence() operatorPrecedence
	initiate(name string)
}

type operatorPrecedence uint8

func (pre operatorPrecedence) merge(op operator) operator {
	if op, ok := op.(operatorPrecedence); ok {
		if op > pre {
			return op
		}
		return pre
	}
	if op == nil {
		return pre
	}
	return op.merge(pre)
}

func (pre operatorPrecedence) precedence() operatorPrecedence {
	return pre
}

func (pre operatorPrecedence) initiate(name string) {}

type infix struct {
	operatorPrecedence
	number       func(a, b float64) (interface{}, error)
	decimal      func(a, b decimal.Decimal) (interface{}, error)
	boolean      func(a, b bool) (interface{}, error)
	text         func(a, b string) (interface{}, error)
	arbitrary    func(a, b interface{}) (interface{}, error)
	shortCircuit func(a interface{}) (interface{}, bool)
	builder      infixBuilder
}

func (op infix) merge(op2 operator) operator {
	switch op2 := op2.(type) {
	case *infix:
		if op.number == nil {
			op.number = op2.number
		}
		if op.decimal == nil {
			op.decimal = op2.decimal
		}
		if op.boolean == nil {
			op.boolean = op2.boolean
		}
		if op.text == nil {
			op.text = op2.text
		}
		if op.arbitrary == nil {
			op.arbitrary = op2.arbitrary
		}
		if op.shortCircuit == nil {
			op.shortCircuit = op2.shortCircuit
		}
	}
	if op2 != nil && op2.precedence() > op.operatorPrecedence {
		op.operatorPrecedence = op2.precedence()
	}
	return &op
}

type directInfix struct {
	operatorPrecedence
	infixBuilder
}

func (op directInfix) merge(op2 operator) operator {
	switch op2 := op2.(type) {
	case operatorPrecedence:
		op.operatorPrecedence = op2
	}
	if op2 != nil && op2.precedence() > op.operatorPrecedence {
		op.operatorPrecedence = op2.precedence()
	}
	return op
}

type extension func(context.Context, *Parser) (Evaluable, error)

type postfix struct {
	operatorPrecedence
	f func(context.Context, *Parser, Evaluable, operatorPrecedence) (Evaluable, error)
}

func (op postfix) merge(op2 operator) operator {
	switch op2 := op2.(type) {
	case postfix:
		if op2.f != nil {
			op.f = op2.f
		}
	}
	if op2 != nil && op2.precedence() > op.operatorPrecedence {
		op.operatorPrecedence = op2.precedence()
	}
	return op
}
