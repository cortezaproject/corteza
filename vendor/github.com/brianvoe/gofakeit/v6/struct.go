package gofakeit

import (
	"errors"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Struct fills in exported fields of a struct with random data
// based on the value of `fake` tag of exported fields.
// Use `fake:"skip"` to explicitly skip an element.
// All built-in types are supported, with templating support
// for string types.
func Struct(v interface{}) error { return structFunc(globalFaker.Rand, v) }

// Struct fills in exported fields of a struct with random data
// based on the value of `fake` tag of exported fields.
// Use `fake:"skip"` to explicitly skip an element.
// All built-in types are supported, with templating support
// for string types.
func (f *Faker) Struct(v interface{}) error { return structFunc(f.Rand, v) }

func structFunc(ra *rand.Rand, v interface{}) error {
	return r(ra, reflect.TypeOf(v), reflect.ValueOf(v), "", 0)
}

func r(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string, size int) error {
	switch t.Kind() {
	case reflect.Ptr:
		return rPointer(ra, t, v, tag, size)
	case reflect.Struct:
		return rStruct(ra, t, v, tag)
	case reflect.String:
		return rString(ra, t, v, tag)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rUint(ra, t, v, tag)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rInt(ra, t, v, tag)
	case reflect.Float32, reflect.Float64:
		return rFloat(ra, t, v, tag)
	case reflect.Bool:
		return rBool(ra, t, v, tag)
	case reflect.Array, reflect.Slice:
		return rSlice(ra, t, v, tag, size)
	}

	return nil
}

func rStruct(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string) error {
	// If tag is set lets try to set the struct values from the tag response
	if tag != "" {
		// Trim the curly on the begining and end
		tag = strings.TrimLeft(tag, "{")
		tag = strings.TrimRight(tag, "}")

		// Check if has params separated by :
		fNameSplit := strings.SplitN(tag, ":", 2)
		fName := ""
		fParams := ""
		if len(fNameSplit) >= 1 {
			fName = fNameSplit[0]
		}
		if len(fNameSplit) >= 2 {
			fParams = fNameSplit[1]
		}

		// Check to see if its a replaceable lookup function
		if info := GetFuncLookup(fName); info != nil {
			// Get parameters, make sure params and the split both have values
			var mapParams *MapParams
			paramsLen := len(info.Params)

			// If just one param and its a string simply just pass it
			if paramsLen == 1 && info.Params[0].Type == "string" {
				if mapParams == nil {
					mapParams = NewMapParams()
				}
				mapParams.Add(info.Params[0].Field, fParams)
			} else if paramsLen > 0 && fParams != "" {
				splitVals := funcLookupSplit(fParams)
				for ii := 0; ii < len(splitVals); ii++ {
					if paramsLen-1 >= ii {
						if mapParams == nil {
							mapParams = NewMapParams()
						}
						if strings.HasPrefix(splitVals[ii], "[") {
							lookupSplits := funcLookupSplit(strings.TrimRight(strings.TrimLeft(splitVals[ii], "["), "]"))
							for _, v := range lookupSplits {
								mapParams.Add(info.Params[ii].Field, v)
							}
						} else {
							mapParams.Add(info.Params[ii].Field, splitVals[ii])
						}
					}
				}
			}

			// Call function
			fValue, err := info.Generate(ra, mapParams, info)
			if err != nil {
				return err
			}

			// Create new element of expected type
			field := reflect.New(reflect.TypeOf(fValue))
			field.Elem().Set(reflect.ValueOf(fValue))

			// Check if element is pointer if so
			// grab the underlyning value before setting
			fieldElem := field.Elem()
			if fieldElem.Kind() == reflect.Ptr {
				v.Set(fieldElem.Elem())
			} else {
				v.Set(fieldElem)
			}

			// If a function is called to set the struct
			// stop from going through sub fields
			return nil
		}
	}

	n := t.NumField()
	for i := 0; i < n; i++ {
		elementT := t.Field(i)
		elementV := v.Field(i)
		fakeTag, ok := elementT.Tag.Lookup("fake")

		// Check whether or not to skip this field
		if ok && fakeTag == "skip" {
			// Do nothing, skip it
			continue
		}

		// Check to make sure you can set it or that its an embeded(anonymous) field
		if elementV.CanSet() || elementT.Anonymous {
			// Check if reflect type is of values we can specifically set
			switch elementT.Type.String() {
			case "time.Time":
				err := rTime(ra, elementT, elementV, fakeTag)
				if err != nil {
					return err
				}
				continue
			}

			// Check if fakesize is set
			size := -1 // Set to -1 to indicate fakesize was not set
			fs, ok := elementT.Tag.Lookup("fakesize")
			if ok {
				var err error
				size, err = strconv.Atoi(fs)
				if err != nil {
					return err
				}
			}
			err := r(ra, elementT.Type, elementV, fakeTag, size)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func rPointer(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string, size int) error {
	elemT := t.Elem()
	if v.IsNil() {
		nv := reflect.New(elemT)
		err := r(ra, elemT, nv.Elem(), tag, size)
		if err != nil {
			return err
		}
		v.Set(nv)
	} else {
		err := r(ra, elemT, v.Elem(), tag, size)
		if err != nil {
			return err
		}
	}

	return nil
}

func rSlice(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string, size int) error {
	// If you cant even set it dont even try
	if !v.CanSet() {
		return errors.New("Cannot set slice")
	}

	// Grab original size to use if needed for sub arrays
	ogSize := size

	// If the value has a cap and is less than the size
	// use that instead of the requested size
	elemCap := v.Cap()
	if elemCap == 0 && size == -1 {
		size = number(ra, 1, 10)
	} else if elemCap != 0 && (size == -1 || elemCap < size) {
		size = elemCap
	}

	// Get the element type
	elemT := t.Elem()

	// If values are already set fill them up, otherwise append
	if v.Len() != 0 {
		// Loop through the elements length and set based upon the index
		for i := 0; i < size; i++ {
			nv := reflect.New(elemT)
			err := r(ra, elemT, nv.Elem(), tag, ogSize)
			if err != nil {
				return err
			}
			v.Index(i).Set(reflect.Indirect(nv))
		}
	} else {
		// Loop through the size and append and set
		for i := 0; i < size; i++ {
			nv := reflect.New(elemT)
			err := r(ra, elemT, nv.Elem(), tag, ogSize)
			if err != nil {
				return err
			}
			v.Set(reflect.Append(reflect.Indirect(v), reflect.Indirect(nv)))
		}
	}

	return nil
}

func rString(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string) error {
	if tag != "" {
		v.SetString(generate(ra, tag))
	} else {
		v.SetString(generate(ra, strings.Repeat("?", number(ra, 4, 10))))
	}

	return nil
}

func rInt(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string) error {
	if tag != "" {
		i, err := strconv.ParseInt(generate(ra, tag), 10, 64)
		if err != nil {
			return err
		}

		v.SetInt(i)
		return nil
	}

	// If no tag or error converting to int, set with random value
	switch t.Kind() {
	case reflect.Int:
		v.SetInt(int64Func(ra))
	case reflect.Int8:
		v.SetInt(int64(int8Func(ra)))
	case reflect.Int16:
		v.SetInt(int64(int16Func(ra)))
	case reflect.Int32:
		v.SetInt(int64(int32Func(ra)))
	case reflect.Int64:
		v.SetInt(int64Func(ra))
	}

	return nil
}

func rUint(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string) error {
	if tag != "" {
		u, err := strconv.ParseUint(generate(ra, tag), 10, 64)
		if err != nil {
			return err
		}

		v.SetUint(u)
		return nil
	}

	// If no tag or error converting to uint, set with random value
	switch t.Kind() {
	case reflect.Uint:
		v.SetUint(uint64Func(ra))
	case reflect.Uint8:
		v.SetUint(uint64(uint8Func(ra)))
	case reflect.Uint16:
		v.SetUint(uint64(uint16Func(ra)))
	case reflect.Uint32:
		v.SetUint(uint64(uint32Func(ra)))
	case reflect.Uint64:
		v.SetUint(uint64Func(ra))
	}

	return nil
}

func rFloat(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string) error {
	if tag != "" {
		f, err := strconv.ParseFloat(generate(ra, tag), 64)
		if err != nil {
			return err
		}

		v.SetFloat(f)
		return nil
	}

	// If no tag or error converting to float, set with random value
	switch t.Kind() {
	case reflect.Float64:
		v.SetFloat(float64Func(ra))
	case reflect.Float32:
		v.SetFloat(float64(float32Func(ra)))
	}

	return nil
}

func rBool(ra *rand.Rand, t reflect.Type, v reflect.Value, tag string) error {
	if tag != "" {
		b, err := strconv.ParseBool(generate(ra, tag))
		if err != nil {
			return err
		}

		v.SetBool(b)
		return nil
	}

	// If no tag or error converting to boolean, set with random value
	v.SetBool(boolFunc(ra))

	return nil
}

// rTime will set a time.Time field the best it can from either the default date tag or from the generate tag
func rTime(ra *rand.Rand, t reflect.StructField, v reflect.Value, tag string) error {
	if tag != "" {
		timeFormat, timeFormatOK := t.Tag.Lookup("format")
		if !timeFormatOK {
			timeFormat = time.RFC3339
		}

		timeFormat = javaDateFormatToGolangDateFormat(timeFormat)

		// Generate time
		timeOutput := generate(ra, tag)

		// If output is larger than format cut the output
		if len(timeOutput) > len(timeFormat) {
			timeOutput = timeOutput[:len(timeFormat)]
		}

		timeStruct, err := time.Parse(timeFormat, timeOutput)
		if err != nil {
			return err
		}

		v.Set(reflect.ValueOf(timeStruct))
		return nil
	}

	v.Set(reflect.ValueOf(date(ra)))
	return nil
}
