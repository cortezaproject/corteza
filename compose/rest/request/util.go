package request

//lint:file-ignore U1000 Ignore unused code, part of request pkg toolset

import (
	"encoding/json"
	"strconv"
	"strings"
)

type (
	ProcedureArgs []ProcedureArg

	ProcedureArg struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
)

func (args ProcedureArgs) GetUint64(name string) uint64 {
	u, _ := strconv.ParseUint(args.Get(name), 10, 64)
	return u
}

func (args ProcedureArgs) Get(name string) string {
	name = strings.ToLower(name)
	for _, arg := range args {
		if strings.ToLower(arg.Name) == name {
			return arg.Value
		}
	}

	return ""
}

func parseMapStringInterface(ss []string) (map[string]interface{}, error) {
	if len(ss) == 0 {
		return nil, nil
	}

	out := make(map[string]interface{})

	return out, json.Unmarshal([]byte(ss[0]), out)
}
