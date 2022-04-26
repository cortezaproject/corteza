package types

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cortezaproject/corteza-server/pkg/locale"
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
		if aux = tt.FindByKey(LocaleKeyChartYAxisLabel.Path); aux != nil {
			c.Config.Reports[i].YAxis["label"] = aux.Msg
		}

		for j, metric := range report.Metrics {
			if metricID, ok := metric["metricID"]; ok {
				mpl := strings.NewReplacer(
					"{{metricID}}", metricID.(string),
				)

				if aux = tt.FindByKey(mpl.Replace(LocaleKeyChartMetricsMetricIDLabel.Path)); aux != nil {
					c.Config.Reports[i].Metrics[j]["label"] = aux.Msg
				}
			}
		}
	}
}

func (c Chart) encodeTranslations() (out locale.ResourceTranslationSet) {
	out = make(locale.ResourceTranslationSet, 0, 12)

	for _, report := range c.Config.Reports {
		if _, ok := report.YAxis["label"]; ok {
			out = append(out, &locale.ResourceTranslation{
				Resource: c.ResourceTranslation(),
				Key:      LocaleKeyChartYAxisLabel.Path,
				Msg:      report.YAxis["label"].(string),
			})
		}

		for _, metric := range report.Metrics {
			if metricID, ok := metric["metricID"]; ok {
				mpl := strings.NewReplacer(
					"{{metricID}}", metricID.(string),
				)

				if _, ok = metric["label"]; ok {
					out = append(out, &locale.ResourceTranslation{
						Resource: c.ResourceTranslation(),
						Key:      mpl.Replace(LocaleKeyChartMetricsMetricIDLabel.Path),
						Msg:      metric["label"].(string),
					})
				}
			}
		}
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
