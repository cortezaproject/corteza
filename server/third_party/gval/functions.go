package gval

import (
	"context"
	"fmt"
	"reflect"
)

type function func(ctx context.Context, arguments ...interface{}) (interface{}, error)

func toFunc(f interface{}) function {
	if f, ok := f.(func(arguments ...interface{}) (interface{}, error)); ok {
		return function(func(ctx context.Context, arguments ...interface{}) (interface{}, error) {
			var v interface{}
			errCh := make(chan error, 1)
			go func() {
				defer func() {
					if recovered := recover(); recovered != nil {
						errCh <- fmt.Errorf("%v", recovered)
					}
				}()
				result, err := f(arguments...)
				v = result
				errCh <- err
			}()

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case err := <-errCh:
				close(errCh)
				return v, err
			}
		})
	}
	if f, ok := f.(func(ctx context.Context, arguments ...interface{}) (interface{}, error)); ok {
		return function(f)
	}

	fun := reflect.ValueOf(f)
	t := fun.Type()
	return func(ctx context.Context, args ...interface{}) (interface{}, error) {
		var v interface{}
		errCh := make(chan error, 1)
		go func() {
			defer func() {
				if recovered := recover(); recovered != nil {
					errCh <- fmt.Errorf("%v", recovered)
				}
			}()
			in, err := createCallArguments(ctx, t, args)
			if err != nil {
				errCh <- err
				return
			}
			out := fun.Call(in)

			r := make([]interface{}, len(out))
			for i, e := range out {
				r[i] = e.Interface()
			}

			err = nil
			errorInterface := reflect.TypeOf((*error)(nil)).Elem()
			if len(r) > 0 && t.Out(len(r)-1).Implements(errorInterface) {
				if r[len(r)-1] != nil {
					err = r[len(r)-1].(error)
				}
				r = r[0 : len(r)-1]
			}

			switch len(r) {
			case 0:
				v = nil
			case 1:
				v = r[0]
			default:
				v = r
			}
			errCh <- err
		}()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errCh:
			close(errCh)
			return v, err
		}
	}
}

func createCallArguments(ctx context.Context, t reflect.Type, args []interface{}) ([]reflect.Value, error) {
	variadic := t.IsVariadic()
	numIn := t.NumIn()

	// if first argument is a context, use the given execution context
	if numIn > 0 {
		thisFun := reflect.ValueOf(createCallArguments)
		thisT := thisFun.Type()
		if t.In(0) == thisT.In(0) {
			args = append([]interface{}{ctx}, args...)
		}
	}

	if (!variadic && len(args) != numIn) || (variadic && len(args) < numIn-1) {
		return nil, fmt.Errorf("invalid number of parameters")
	}

	in := make([]reflect.Value, len(args))
	var inType reflect.Type
	for i, arg := range args {
		if !variadic || i < numIn-1 {
			inType = t.In(i)
		} else if i == numIn-1 {
			inType = t.In(numIn - 1).Elem()
		}
		argVal := reflect.ValueOf(arg)
		if arg == nil {
			argVal = reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem())
		} else if !argVal.Type().AssignableTo(inType) {
			return nil, fmt.Errorf("expected type %s for parameter %d but got %T",
				inType.String(), i, arg)
		}
		in[i] = argVal
	}
	return in, nil
}
