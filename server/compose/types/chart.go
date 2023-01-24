package types

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/cortezaproject/corteza/server/pkg/sql"
	"github.com/spf13/cast"

	"github.com/cortezaproject/corteza/server/pkg/filter"
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
		NoAnimation bool                 `json:"noAnimation,omitempty"`
	}

	ChartConfigReport struct {
		ReportID   uint64                   `json:"reportID,string,omitempty"`
		Filter     string                   `json:"filter"`
		ModuleID   uint64                   `json:"moduleID,string,omitempty"`
		Metrics    []map[string]interface{} `json:"metrics,omitempty"`
		Dimensions []map[string]interface{} `json:"dimensions,omitempty"`
		YAxis      map[string]interface{}   `json:"yAxis,omitempty"`
		Legend     map[string]interface{}   `json:"legend,omitempty"`
		Tooltip    map[string]interface{}   `json:"tooltip,omitempty"`
		Offset     map[string]interface{}   `json:"offset,omitempty"`
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
		report.WalkMetrics(func(metricID string, metric map[string]interface{}) {
			mpl := strings.NewReplacer("{{metricID}}", metricID)

			aux = tt.FindByKey(mpl.Replace(LocaleKeyChartMetricsMetricIDLabel.Path))
			if aux == nil {
				return
			}

			metric["label"] = aux.Msg
		})

		// apply translated labels for each dimension/step
		report.WalkDimensionSteps(func(dimensionID, stepID string, step map[string]interface{}) {
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
		report.WalkMetrics(func(metricID string, m map[string]interface{}) {
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
		report.WalkDimensionSteps(func(dimID, stepID string, step map[string]interface{}) {
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

func (cc *ChartConfig) Scan(src any) error          { return sql.ParseJSON(src, cc) }
func (cc ChartConfig) Value() (driver.Value, error) { return json.Marshal(cc) }

func (r *ChartConfigReport) WalkMetrics(fn func(string, map[string]interface{})) {
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

func (r *ChartConfigReport) WalkDimensionSteps(fn func(string, string, map[string]interface{})) {
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

func (c *ChartConfig) GenerateIDs(nextID func() uint64) {
	// Ensure chart report IDs
	for r := range c.Reports {
		c.Reports[r].ReportID = nextID()

		// Ensure chart report metric IDs
		for m := range c.Reports[r].Metrics {
			met := c.Reports[r].Metrics[m]
			if _, has := met["metricID"]; has {
				continue
			}

			met["metricID"] = strconv.FormatUint(nextID(), 10)
		}

		for d := range c.Reports[r].Dimensions {
			dim := c.Reports[r].Dimensions[d]

			if _, has := dim["dimensionID"]; !has {
				dim["dimensionID"] = strconv.FormatUint(nextID(), 10)
			}

			meta, is := dim["meta"].(map[string]interface{})
			if !is {
				// no meta, no steps
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
				_, has := steps[s]["stepID"]
				if has {
					continue
				}

				steps[s]["stepID"] = strconv.FormatUint(nextID(), 10)
			}

			meta["steps"] = steps
			dim["meta"] = meta
		}
	}

}
