package options

import (
	"os"
	"reflect"
	"time"

	"github.com/spf13/cast"
)

func fill(opt interface{}) {
	v := reflect.ValueOf(opt)
	if v.Kind() != reflect.Ptr {
		panic("expecting a pointer, not a value")
	}

	if v.IsNil() {
		panic("nil pointer passed")
	}

	v = v.Elem()

	length := v.NumField()
	for i := 0; i < length; i++ {
		f := v.Field(i)
		t := v.Type().Field(i)
		if tag := t.Tag.Get("env"); tag != "" {
			if !f.CanSet() {
				panic("unexpected pointer for field " + t.Name)
			}

			if f.Type() == reflect.TypeOf(time.Duration(1)) {
				v.FieldByName(t.Name).SetInt(int64(EnvDuration(tag, time.Duration(f.Int()))))
				continue
			}

			if f.Kind() == reflect.String {
				v.FieldByName(t.Name).SetString(EnvString(tag, f.String()))
				continue
			}

			if f.Kind() == reflect.Bool {
				v.FieldByName(t.Name).SetBool(EnvBool(tag, f.Bool()))
				continue
			}

			if f.Kind() == reflect.Int {
				v.FieldByName(t.Name).SetInt(int64(EnvInt(tag, int(f.Int()))))
				continue
			}

			if f.Kind() == reflect.Float32 {
				v.FieldByName(t.Name).SetFloat(float64(EnvFloat32(tag, float32(f.Float()))))
				continue
			}

			panic("unsupported type/kind for field " + t.Name)

		}
	}
}

func EnvString(key string, def string) string {
	if val, has := os.LookupEnv(key); has {
		return val
	}
	return def
}

func EnvBool(key string, def bool) bool {
	if val, has := os.LookupEnv(key); has {
		if b, err := cast.ToBoolE(val); err == nil {
			return b
		}
	}
	return def
}

func EnvInt(key string, def int) int {
	if val, has := os.LookupEnv(key); has {
		if i, err := cast.ToIntE(val); err == nil {
			return i
		}
	}
	return def
}

func EnvFloat32(key string, def float32) float32 {
	if val, has := os.LookupEnv(key); has {
		if i, err := cast.ToFloat32E(val); err == nil {
			return i
		}
	}
	return def
}

func EnvDuration(key string, def time.Duration) time.Duration {
	if val, has := os.LookupEnv(key); has {
		if d, err := cast.ToDurationE(val); err == nil {
			return d
		}
	}
	return def
}
