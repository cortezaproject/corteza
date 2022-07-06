package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/spf13/cast"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/pkg/errors"
)

type (
	Chart struct {
		ID     uint64      `json:"chartID,string"`
		Handle string      `json:"handle"`
		Name   string      `json:"name"`
		Config ChartConfig `json:"config"`

		Labels map[string]string `json:"labels,omitempty"`

		NamespaceID uint64 `json:"namespaceID,string"`

		CreatedAt time.Time  `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `json:"deletedAt,omitempty"`
	}

	ChartConfig struct {
		Reports     []*ChartConfigReport `json:"reports,omitempty"`
		ColorScheme string               `json:"colorScheme,omitempty"`
	}

	ChartConfigReport struct {
		ReportID   uint64                   `json:"reportID,string,omitempty"`
		Filter     string                   `json:"filter"`
		ModuleID   uint64                   `json:"moduleID,string,omitempty"`
		Metrics    []map[string]interface{} `json:"metrics,omitempty"`
		Dimensions []map[string]interface{} `json:"dimensions,omitempty"`
		YAxis      map[string]interface{}   `json:"yAxis,omitempty"`
		Renderer   struct {
			Version string `json:"version,omitempty" `
		} `json:"renderer,omitempty"`
	}

	ChartFilter struct {
		NamespaceID uint64   `json:"namespaceID,string"`
		ChartID     []uint64 `json:"chartID"`
		Handle      string   `json:"handle"`
		Name        string   `json:"name"`
		Query       string   `json:"query"`

		LabeledIDs []uint64          `json:"-"`
		Labels     map[string]string `json:"labels,omitempty"`

		Deleted filter.State `json:"deleted"`

		// Check fn is called by store backend for each resource found function can
		// modify the resource and return false if store should not return it
		//
		// Store then loads additional resources to satisfy the paging parameters
		Check func(*Chart) (bool, error) `json:"-"`

		// Standard helpers for paging and sorting
		filter.Sorting
		filter.Paging
	}
)

func (c Chart) decodeTranslations(tt locale.ResourceTranslationIndex) {
	var aux *locale.ResourceTranslation

	for i, report := range c.Config.Reports {
		if report == nil {
			continue
		}

		// apply translated label for YAxis
		if aux = tt.FindByKey(LocaleKeyChartYAxisLabel.Path); aux != nil {
			if c.Config.Reports[i].YAxis == nil {
				c.Config.Reports[i].YAxis = make(map[string]interface{})
			}

			c.Config.Reports[i].YAxis["label"] = aux.Msg
		}

		// apply translated labels for metrics
		report.walkMetrics(func(metricID string, metric map[string]interface{}) {
			mpl := strings.NewReplacer("{{metricID}}", metricID)

			aux = tt.FindByKey(mpl.Replace(LocaleKeyChartMetricsMetricIDLabel.Path))
			if aux == nil {
				return
			}

			metric["label"] = aux.Msg
		})

		// apply translated labels for each dimension/step
		report.walkDimensionSteps(func(dimensionID, stepID string, step map[string]interface{}) {
			mpl := strings.NewReplacer(
				"{{dimensionID}}", dimensionID,
				"{{stepID}}", stepID,
			)

			aux = tt.FindByKey(mpl.Replace(LocaleKeyChartDimensionsDimensionIDMetaStepsStepIDLabel.Path))
			if aux == nil {
				return
			}

			step["label"] = aux.Msg
		})
	}
}

func (c Chart) encodeTranslations() (out locale.ResourceTranslationSet) {
	out = make(locale.ResourceTranslationSet, 0, 12)

	for _, report := range c.Config.Reports {
		// collect labels from chart config: YAxis
		if _, ok := report.YAxis["label"]; ok {
			out = append(out, &locale.ResourceTranslation{
				Resource: c.ResourceTranslation(),
				Key:      LocaleKeyChartYAxisLabel.Path,
				Msg:      cast.ToString(report.YAxis["label"]),
			})
		}

		// collect labels from chart config: metrics
		report.walkMetrics(func(metricID string, m map[string]interface{}) {
			mpl := strings.NewReplacer(
				"{{metricID}}", metricID,
			)

			out = append(out, &locale.ResourceTranslation{
				Resource: c.ResourceTranslation(),
				Key:      mpl.Replace(LocaleKeyChartMetricsMetricIDLabel.Path),
				Msg:      cast.ToString(m["label"]),
			})
		})

		// collect labels from chart config: dimensions/steps
		report.walkDimensionSteps(func(dimID, stepID string, step map[string]interface{}) {
			mpl := strings.NewReplacer(
				"{{dimensionID}}", dimID,
				"{{stepID}}", stepID,
			)

			out = append(out, &locale.ResourceTranslation{
				Resource: c.ResourceTranslation(),
				Key:      mpl.Replace(LocaleKeyChartDimensionsDimensionIDMetaStepsStepIDLabel.Path),
				Msg:      cast.ToString(step["label"]),
			})
		})
	}

	return
}

// FindByHandle finds chart by it's handle
func (set ChartSet) FindByHandle(handle string) *Chart {
	for i := range set {
		if set[i].Handle == handle {
			return set[i]
		}
	}

	return nil
}

func (cc *ChartConfig) Scan(value interface{}) error {
	//lint:ignore S1034 This typecast is intentional, we need to get []byte out of a []uint8
	switch value.(type) {
	case nil:
		*cc = ChartConfig{}
	case []uint8:
		b := value.([]byte)
		if err := json.Unmarshal(b, cc); err != nil {
			return errors.Wrapf(err, "cannot scan '%v' into ChartConfig", string(b))
		}
	}

	return nil
}

func (cc ChartConfig) Value() (driver.Value, error) {
	return json.Marshal(cc)
}

func (r *ChartConfigReport) walkMetrics(fn func(string, map[string]interface{})) {
	for m := range r.Metrics {
		metricID, ok := r.Metrics[m]["metricID"]
		if !ok {
			continue
		}

		if len(r.Metrics[m]) == 0 {
			// avoid problems with nil maps
			r.Metrics[m] = make(map[string]interface{})
		}

		fn(metricID.(string), r.Metrics[m])
	}
}

func (r *ChartConfigReport) walkDimensionSteps(fn func(string, string, map[string]interface{})) {
	for d := range r.Dimensions {
		dimensionID, ok := r.Dimensions[d]["dimensionID"]
		if !ok {
			continue
		}

		meta, is := r.Dimensions[d]["meta"].(map[string]interface{})
		if !is {
			continue
		}

		var steps []map[string]interface{}

		switch aux := meta["steps"].(type) {
		case []interface{}:
			for _, i := range aux {
				if kv, is := i.(map[string]interface{}); is {
					steps = append(steps, kv)
				}
			}
		case []map[string]interface{}:
			steps = aux
		}

		for s := range steps {
			stepID, has := steps[s]["stepID"]
			if !has {
				return
			}

			fn(dimensionID.(string), stepID.(string), steps[s])
		}
	}
}

// Sets new value into nested map/slice struct
//
// This is a utility function, should be moved
// somewhere under pkg/.
func set(d interface{}, value interface{}, path ...string) {
	index := func(v reflect.Value, idx string) reflect.Value {
		if i, err := strconv.Atoi(idx); err == nil {
			return v.Index(i)
		}
		return v.FieldByName(idx)
	}

	v := reflect.ValueOf(d)
	for _, s := range path {
		v = index(v, s)
	}

	v.Set(reflect.ValueOf(value))
}
