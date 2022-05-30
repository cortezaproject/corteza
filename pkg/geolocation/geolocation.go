package geolocation

import (
	"database/sql/driver"
	"encoding/json"
)

type (
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	}

	Properties struct {
		Name string `json:"name"`
	}

	Full struct {
		Geometry   Geometry   `json:"geometry"`
		Properties Properties `json:"properties"`
	}
)

func Parse(ss []string) (m Full, err error) {
	if len(ss) == 0 {
		return
	}

	err = json.Unmarshal([]byte(ss[0]), &m)
	return
}

func (set *Full) Scan(src interface{}) error {
	if data, ok := src.([]byte); ok {
		return json.Unmarshal(data, set)
	}
	return nil
}

func (set Full) Value() (driver.Value, error) {
	return json.Marshal(set)
}
