package cast2

import (
	"encoding/json"
	"fmt"
	"github.com/modern-go/reflect2"
)

// Meta casts a value to a map[string]any
func Meta(in any, out *map[string]any) (err error) {
	if reflect2.IsNil(in) {
		*out = nil
		return nil
	}

	switch aux := in.(type) {
	case []byte:
		err = json.Unmarshal(aux, out)
	case string:
		err = json.Unmarshal([]byte(aux), out)
	case map[string]any:
		*out = aux
	default:
		err = fmt.Errorf("unsupported type: %T", in)
	}

	if err == nil {
		return
	}

	return fmt.Errorf("can not cast to Meta: %w", err)
}
