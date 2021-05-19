package httperr

// ErrorTranslatorFunc is a function that translates errors
// before they are returned by the http client.
type ErrorTranslatorFunc func(err error) error

type errorTranslatorFuncHolder struct {
	F ErrorTranslatorFunc
}

// RegisterErrorTranslator adds a new error translation function to the
// error translation stack. Your function will be called with each
// error before writing the response.
//
// Note: this function is not go-routine safe. You should probably call
// it from an init() function or similar.
//
// Example:
//
//   init() {
//     RegisterErrorTranslator(func(err error) error {
//       if (err == context.DeadlineExceeded) {
//         return Error{
//           PrivateError: err,
//           StatusCode:   http.StatusRequestTimeout,
//         }
//       }
//     })
//   }
//
func RegisterErrorTranslator(f ErrorTranslatorFunc) {
	holder := &errorTranslatorFuncHolder{F: f}
	errorTranslators = append([]*errorTranslatorFuncHolder{holder}, errorTranslators...)
}

var errorTranslators []*errorTranslatorFuncHolder

// TranslateError invokes all the registered ErrorTranslatorFunc functions on
// err and returns the result.
func TranslateError(err error) error {
	for _, et := range errorTranslators {
		err = et.F(err)
	}
	return err
}
