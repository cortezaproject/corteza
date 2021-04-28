package httperr

import (
	"net/http"
)

// Response is an alias for http.Response that implements
// the error interface. Example:
//
//   resp, err := http.Get("http://www.example.com")
//   if err != nil {
//   	return err
//   }
//   if resp.StatusCode != http.StatusOK {
//   	return httperr.Response(*resp)
//   }
//   // ...
//
type Response http.Response

func (re Response) Error() string {
	msg := re.Header.Get("X-Error-Message")
	if msg != "" {
		msg = ": " + msg
	}
	return re.Status + msg
}
